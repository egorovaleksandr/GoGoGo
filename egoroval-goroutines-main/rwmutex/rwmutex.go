//go:build !solution

package rwmutex

type RWMutex struct {
	num int
	r   chan bool
	w   chan bool
}

func New() *RWMutex {
	res := &RWMutex{0, make(chan bool, 1), make(chan bool, 1)}
	res.r <- false
	res.w <- false
	return res
}

func (rw *RWMutex) RLock() {
	<-rw.r
	if rw.num == 0 {
		<-rw.w
	}
	rw.num++
	rw.r <- false
}

func (rw *RWMutex) RUnlock() {
	<-rw.r
	if rw.num == 1 {
		rw.w <- false
	}
	rw.num--
	rw.r <- false
}

func (rw *RWMutex) Lock() {
	<-rw.w
}

func (rw *RWMutex) Unlock() {
	rw.w <- false
}
