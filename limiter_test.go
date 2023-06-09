package limiter

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"time"
)

var limits = []Limit{
	500 * 1024,  // 500KB
	1024 * 1024, // 1MB
}

var sizes = []int64{
	1024 * 1024,     // 1MB
	1024 * 1024 * 2, // 2MB
	1024 * 1024 * 3, // 3MB
}

func TestLimitRead(t *testing.T) {
	for _, size := range sizes {
		for _, limit := range limits {
			b := make([]byte, size)
			rand.Read(b)
			var buf bytes.Buffer
			limiter := New(limit)
			start := time.Now()
			n, err := io.Copy(&buf, limiter.Reader(bytes.NewReader(b)))
			if err != nil {
				t.Fatal(err)
			}
			if n != size {
				t.Fatalf("%.f/sec expect length %d; got %d", limit, size, n)
			}
			if !bytes.Equal(b, buf.Bytes()) {
				t.Fatal("bytes not equal")
			}
			elapsed := int64(time.Since(start).Truncate(time.Second).Seconds())
			if s := size/int64(limit) - 1; s != elapsed {
				t.Fatalf("%.f/sec %d expect elapsed time %ds; got %ds", limit, size, s, n)
			}
		}
	}
}

func TestLimitWrite(t *testing.T) {
	for _, size := range sizes {
		for _, limit := range limits {
			b := make([]byte, size)
			rand.Read(b)
			var buf bytes.Buffer
			limiter := New(limit)
			start := time.Now()
			n, err := io.Copy(limiter.Writer(&buf), bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			if n != size {
				t.Fatalf("%.f/sec expect length %d; got %d", limit, size, n)
			}
			if !bytes.Equal(b, buf.Bytes()) {
				t.Fatal("bytes not equal")
			}
			elapsed := int64(time.Since(start).Truncate(time.Second).Seconds())
			if s := size/int64(limit) - 1; s != elapsed {
				t.Fatalf("%.f/sec %d expect elapsed time %ds; got %ds", limit, size, s, n)
			}
		}
	}
}
