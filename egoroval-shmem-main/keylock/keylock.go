//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct {
	locked map[string]chan int
	lmut   sync.Mutex
}

func New() *KeyLock {
	return &KeyLock{locked: make(map[string]chan int)}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	canceled = false
	keys_copy := make([]string, len(keys))
	copy(keys_copy, keys)
	sort.Strings(keys_copy)
	var infch = make([]chan int, 0)
	unlock = func() {
		for _, lock := range infch {
			close(lock)
		}
		l.lmut.Lock()
		for _, key := range keys_copy {
			delete(l.locked, key)
		}
		l.lmut.Unlock()
	}
	for _, key := range keys_copy {
		for !canceled {
			l.lmut.Lock()
			loth, larl := l.locked[key]
			lock := make(chan int)
			if !larl {
				l.locked[key] = lock
			}
			l.lmut.Unlock()
			if !larl {
				infch = append(infch, lock)
				break
			}
			select {
			case <-loth:
			case <-cancel:
				canceled = true
				for _, lock := range infch {
					close(lock)
				}
				l.lmut.Lock()
				for _, key := range keys_copy {
					delete(l.locked, key)
				}
				l.lmut.Unlock()
				return
			}
		}
	}
	return
}
