//go:build !solution

package waitgroup

type WaitGroup struct {
	num int
	f   chan struct{}
	s   chan struct{}
}

func New() *WaitGroup {
	r := &WaitGroup{0, make(chan struct{}, 1), make(chan struct{}, 1)}
	r.f <- struct{}{}
	r.s <- struct{}{}
	return r
}

func (wg *WaitGroup) Add(delta int) {
	<-wg.f
	if wg.num == 0 {
		<-wg.s
	}
	wg.num += delta
	if wg.num < 0 {
		panic("negative WaitGroup counter")
	}
	wg.f <- struct{}{}
}

func (wg *WaitGroup) Done() {
	<-wg.f
	if wg.num <= 0 {
		panic("negative WaitGroup counter")
	}
	wg.num--
	if wg.num == 0 {
		wg.s <- struct{}{}
	}
	wg.f <- struct{}{}
}

func (wg *WaitGroup) Wait() {
	<-wg.s
	wg.s <- struct{}{}
}
