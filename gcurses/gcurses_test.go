package gcurses

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	writer := New()
	buf := new(bytes.Buffer)
	writer.Writer = buf

	writer.Start()

	for i := 0; i < 10; i++ {
		fmt.Fprintf(writer, "%d\n", i)
		time.Sleep(time.Millisecond * 10)
	}

	writer.Stop()

	ret := fmt.Sprintf("%d", buf)
	if ret != "&{[48 10 27 91 48 65 27 91 50 75 13 49 10 27 91 48 65 27 91 50 75 13 50 10 27 91 48 65 27 91 50 75 13 51 10 27 91 48 65 27 91 50 75 13 52 10 27 91 48 65 27 91 50 75 13 53 10 27 91 48 65 27 91 50 75 13 54 10 27 91 48 65 27 91 50 75 13 55 10 27 91 48 65 27 91 50 75 13 56 10 27 91 48 65 27 91 50 75 13 57 10] 0 [0 0 0 0] [48 10 27 91 48 65 27 91 50 75 13 49 10 27 91 48 65 27 91 50 75 13 50 10 27 91 48 65 27 91 50 75 13 51 10 27 91 48 65 27 91 50 75 13 52 10 27 91 48 65 27 91 50 75 13 53 10 27 91 48 65 0 0 0] 0}" {
		t.Fatalf("Test failed")
	}
}
