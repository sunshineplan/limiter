package limiter

import (
	"context"
	"io"
)

var _ io.Reader = &reader{}

type reader struct {
	*Limiter
	r   io.Reader
	ctx context.Context
}

func (r *reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if (r.BufferSize() == 0 && len(p) <= r.lim.Burst()) || len(p) <= r.BufferSize() {
		n, err := r.r.Read(p)
		if err := r.waitN(r.ctx, n); err != nil {
			return n, err
		}
		return n, err
	}
	size := r.BufferSize()
	if size == 0 {
		size = 32 * 1024
	}
	if size > r.lim.Burst() {
		size = r.lim.Burst()
	}
	var read int
	for i := 0; i < len(p); i += size {
		end := i + size
		if end > len(p) {
			end = len(p)
		}
		n, err := r.r.Read(p[i:end])
		read += n
		if err := r.waitN(r.ctx, n); err != nil {
			return read, err
		}
		if err != nil {
			return read, err
		}
	}
	return read, nil
}
