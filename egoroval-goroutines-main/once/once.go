//go:build !solution

package once

type Once struct {
	f chan int
	s chan int
}

func New() *Once {
	return &Once{make(chan int, 1), make(chan int, 1)}
}

func (o *Once) Do(f func()) {
	select {
	case o.f <- 0:
		defer func() { o.s <- 11 }()
		f()
	case <-o.s:
		defer func() { o.s <- 11 }()
	}
}
