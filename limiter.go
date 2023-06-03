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

func New(limit int64, bufferSize int) *Limiter {
	limiter := &Limiter{lim: rate.NewLimiter(rate.Limit(limit), int(limit))}
	limiter.bufferSize.Store(int64(bufferSize))
	return limiter
}

func (lim *Limiter) Limit() int64 {
	return int64(lim.lim.Limit())
}

func (lim *Limiter) Burst() int {
	return lim.lim.Burst()
}

func (lim *Limiter) BufferSize() int {
	return int(lim.bufferSize.Load())
}

func (lim *Limiter) SetLimit(newLimit int64) {
	lim.lim.SetLimit(rate.Limit(newLimit))
}

func (lim *Limiter) SetBurst(newBurst int) {
	lim.lim.SetBurst(newBurst)
}

func (lim *Limiter) SetBufferSize(newBufferSize int) {
	lim.bufferSize.Store(int64(newBufferSize))
}

func (lim *Limiter) waitN(ctx context.Context, n int) error {
	return lim.lim.WaitN(ctx, n)
}

func (lim *Limiter) Wrtier(w io.Writer) io.Writer {
	return lim.WrtierWithContext(context.Background(), w)
}

func (lim *Limiter) WrtierWithContext(ctx context.Context, w io.Writer) io.Writer {
	return &writer{lim, w, ctx}
}

func (lim *Limiter) Reader(r io.Reader) io.Reader {
	return lim.ReaderWithContext(context.Background(), r)
}

func (lim *Limiter) ReaderWithContext(ctx context.Context, r io.Reader) io.Reader {
	return &reader{lim, r, ctx}
}
