package limiter

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	"golang.org/x/time/rate"
)

type Limiter struct {
	mu         sync.Mutex
	lim        *rate.Limiter
	bufferSize atomic.Int64
}

// New creates a new Limiter with the given rate limit and buffer size.
func New(limit int64, bufferSize int) *Limiter {
	limiter := &Limiter{lim: rate.NewLimiter(rate.Limit(limit), int(limit))}
	limiter.bufferSize.Store(int64(bufferSize))
	return limiter
}

// Limit returns the current rate limit.
func (lim *Limiter) Limit() int64 {
	return int64(lim.lim.Limit())
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
func (lim *Limiter) SetLimit(newLimit int64) {
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
