// simple caching with GC
package cache

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

type Item struct {
	Object     interface{}
	Expiration time.Time
}

func (item *Item) isExpired() bool {
	if item.Expiration.IsZero() {
		return false
	}
	return item.Expiration.Before(time.Now())
}

type Cache struct {
	Expiration time.Duration
	items      map[string]*Item
}

func (cache *Cache) String() string {
	var str string
	var keys []string
	for k := range cache.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str += k + "\t" + fmt.Sprintf("%v", cache.items[k].Object) + "\n"
	}

	return str
}

func (cache *Cache) Set(key string, value interface{}) {
	cache.items[key] = &Item{
		Object:     value,
		Expiration: time.Now().Add(cache.Expiration),
	}
}

func (cache *Cache) Get(key string) interface{} {
	item, ok := cache.items[key]
	if !ok || item.isExpired() {
		return nil
	}
	return item.Object
}

func (cache *Cache) Delete(key string) {
	delete(cache.items, key)
}

func (cache *Cache) Add(key string, value interface{}) bool {
	item := cache.Get(key)
	if item != nil {
		return false
	}
	cache.Set(key, value)
	return true
}

func (cache *Cache) Update(key string, value interface{}) bool {
	item := cache.Get(key)
	if item == nil {
		return false
	}
	cache.Set(key, value)
	return true
}

func (cache *Cache) DeleteExpired() {
	for k, v := range cache.items {
		if v.isExpired() {
			cache.Delete(k)
		}
	}
}

func (cache *Cache) DeleteExpiredWithFunc(fn func(key string, value interface{})) {
	for k, v := range cache.items {
		if v.isExpired() {
			fn(k, cache.items[k].Object)
			cache.Delete(k)
		}
	}
}

func (cache *Cache) DeleteAllWithFunc(fn func(key string, value interface{})) {
	for k := range cache.items {
		fn(k, cache.items[k].Object)
		cache.Delete(k)
	}
}

func (cache *Cache) Size() int {
	n := len(cache.items)
	return n
}

func (cache *Cache) Clear() {
	cache.items = map[string]*Item{}
}

func cleaner(cache *Cache, interval time.Duration) {
	ticker := time.Tick(interval)
	for {
		select {
		case <-ticker:
			cache.DeleteExpired()
		}
	}
}

func cleanerWithFunc(cache *Cache, interval time.Duration, fn func(key string, value interface{})) {
	defer cache.DeleteAllWithFunc(fn)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ticker := time.Tick(interval)
	for {
		select {
		case <-ticker:
			cache.DeleteExpiredWithFunc(fn)
		case <-c:
			cache.DeleteAllWithFunc(fn)
			os.Exit(1)
		}
	}
}

func New(expirationTime, cleanupInterval time.Duration) *Cache {
	items := make(map[string]*Item)
	if expirationTime == 0 {
		expirationTime = -1
	}
	cache := &Cache{
		Expiration: expirationTime,
		items:      items,
	}
	go cleaner(cache, cleanupInterval)

	return cache
}

func New2(expirationTime, cleanupInterval time.Duration, fn func(key string, value interface{})) *Cache {
	items := make(map[string]*Item)
	if expirationTime == 0 {
		expirationTime = -1
	}
	cache := &Cache{
		Expiration: expirationTime,
		items:      items,
	}
	go cleanerWithFunc(cache, cleanupInterval, fn)

	return cache
}
