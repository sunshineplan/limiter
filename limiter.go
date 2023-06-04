package limiter

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// Limit defines the maximum transfer speed of data.
// The Limit is represented as a rate limit per second.
// A zero Limit means no events are allowed.
type Limit rate.Limit

// Inf is the infinite rate limit; it allows all events (even if burst is zero).
const Inf = Limit(rate.Inf)

// Every converts a minimum time interval between events to a Limit.
func Every(interval time.Duration) Limit {
	return Limit(rate.Every(interval))
}

type Limiter struct {
	mu         sync.Mutex
	lim        *rate.Limiter
	bufferSize atomic.Int64
}

// New creates a new Limiter with the given rate limit and buffer size.
func New(limit Limit, bufferSize int) *Limiter {
	var b int
	if limit == Inf {
		b = 0
	} else {
		b = int(limit)
	}
	limiter := &Limiter{lim: rate.NewLimiter(rate.Limit(limit), b)}
	limiter.bufferSize.Store(int64(bufferSize))
	return limiter
}

// Limit returns the current rate limit.
func (lim *Limiter) Limit() Limit {
	return Limit(lim.lim.Limit())
}

// Burst returns the current burst size.
func (lim *Limiter) Burst() int {
	return lim.lim.Burst()
}

// BufferSize returns the current buffer size.
func (lim *Limiter) BufferSize() int {
	return int(lim.bufferSize.Load())
}

// SetLimit sets a new rate limit.
func (lim *Limiter) SetLimit(newLimit Limit) {
	lim.lim.SetLimit(rate.Limit(newLimit))
}

// SetBurst sets a new burst size.
func (lim *Limiter) SetBurst(newBurst int) {
	lim.lim.SetBurst(newBurst)
}

// SetBufferSize sets a new buffer size.
func (lim *Limiter) SetBufferSize(newBufferSize int) {
	lim.bufferSize.Store(int64(newBufferSize))
}

// waitN waits for availability of n tokens.
func (lim *Limiter) waitN(ctx context.Context, n int) error {
	return lim.lim.WaitN(ctx, n)
}

// Writer returns a writer with rate limiting.
func (lim *Limiter) Writer(w io.Writer) io.Writer {
	return lim.WriterWithContext(context.Background(), w)
}

// WriterWithContext returns a writer with rate limiting and context.
func (lim *Limiter) WriterWithContext(ctx context.Context, w io.Writer) io.Writer {
	return &writer{lim, w, ctx}
}

// Reader returns a reader with rate limiting.
func (lim *Limiter) Reader(r io.Reader) io.Reader {
	return lim.ReaderWithContext(context.Background(), r)
}

// ReaderWithContext returns a reader with rate limiting and context.
func (lim *Limiter) ReaderWithContext(ctx context.Context, r io.Reader) io.Reader {
	return &reader{lim, r, ctx}
}
