//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

type Limiter struct {
	f  int
	n  int
	ts []time.Timer
	t  time.Duration
	c  chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

func NewLimiter(n int, t time.Duration) *Limiter {
	r := &Limiter{0, n, make([]time.Timer, n), t, make(chan struct{}, 1)}
	for i := range r.ts {
		r.ts[i] = *time.NewTimer(0)
	}
	r.c <- struct{}{}
	return r
}

func (l *Limiter) Acquire(ctx context.Context) error {
	if l.f == 0 {
		num := 0
		for {
			num %= l.n
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-l.ts[num].C:
					l.ts[num] = *time.NewTimer(l.t)
					return nil
				default:
					num++
					continue
				}
			}
		}
	}
	return ErrStopped
}

func (l *Limiter) Stop() {
	defer func() { l.f++ }()
	select {
	case <-l.c:
		close(l.c)
	default:
	}
}
