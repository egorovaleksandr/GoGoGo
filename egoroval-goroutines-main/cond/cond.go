//go:build !solution

package cond

type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L Locker
	c chan chan bool
}

func New(l Locker) *Cond {
	return &Cond{L: l, c: make(chan chan bool, 10000)}
}

func (c *Cond) Wait() {
	ch := make(chan bool, 1)
	c.c <- ch
	c.L.Unlock()
	<-ch
	c.L.Lock()
}

func (c *Cond) Signal() {
	select {
	case ch := <-c.c:
		ch <- false
	default:
	}
}

func (c *Cond) Broadcast() {
	for {
		select {
		case ch := <-c.c:
			ch <- false
		default:
			return
		}
	}
}
