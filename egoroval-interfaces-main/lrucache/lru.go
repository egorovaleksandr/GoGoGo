//go:build !solution

package lrucache

import (
	"container/list"
)

type P struct {
	f int
	s int
}

type LRUCache struct {
	cap_ int
	l_   *list.List
	val_ map[int]*list.Element
}

func (lru *LRUCache) Get(key int) (int, bool) {
	if lru.val_ == nil {
		return 0, false
	}
	v, ok := lru.val_[key]
	if !ok {
		return 0, false
	}
	lru.l_.MoveToFront(v)
	return v.Value.(P).s, true
}

func (lru *LRUCache) Set(key, value int) {
	v, ok := lru.val_[key]
	if !ok {
		if lru.cap_ == 0 {
			return
		}
		if lru.l_.Len() == lru.cap_ {
			delete(lru.val_, lru.l_.Back().Value.(P).f)
			lru.l_.Remove(lru.l_.Back())
		}
		lru.val_[key] = lru.l_.PushFront(P{key, value})
		return
	}
	lru.val_[key].Value = P{v.Value.(P).f, value}
	lru.l_.MoveToFront(v)
}

func (lru *LRUCache) Range(f func(key, value int) bool) {
	for v := lru.l_.Back(); v != nil; v = v.Prev() {
		key := v.Value.(P).f
		value := v.Value.(P).s
		if !f(key, value) {
			break
		}
	}
}

func (lru *LRUCache) Clear() {
	lru.l_ = list.New()
	lru.val_ = make(map[int]*list.Element, lru.cap_)
}

func New(cap int) LRUCache {
	return LRUCache{cap, list.New(), make(map[int]*list.Element, cap)}
}
