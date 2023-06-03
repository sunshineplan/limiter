package limiter

import (
	"context"
	"io"
)

var _ io.Writer = &writer{}

type writer struct {
	*Limiter
	w   io.Writer
	ctx context.Context
}

func (w *writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if (w.BufferSize() == 0 && len(p) <= w.lim.Burst()) || len(p) <= w.BufferSize() {
		n, err := w.w.Write(p)
		if err := w.waitN(w.ctx, n); err != nil {
			return n, err
		}
		return n, err
	}
	size := w.BufferSize()
	if size == 0 {
		size = 32 * 1024
	}
	if size > w.lim.Burst() {
		size = w.lim.Burst()
	}
	var written int
	for i := 0; i < len(p); i += size {
		end := i + size
		if end > len(p) {
			end = len(p)
		}
		n, err := w.w.Write(p[i:end])
		written += n
		if err := w.waitN(w.ctx, n); err != nil {
			return written, err
		}
		if err != nil {
			return written, err
		}
	}
	return written, nil
}
