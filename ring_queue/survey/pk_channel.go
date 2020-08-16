package main

import "sync"

type Channel struct {
	cap int
	ch  chan interface{}
	cnt int
	rw  *sync.RWMutex
}

func NewChannel(cap int) *Channel {
	return &Channel{cap: cap, ch: make(chan interface{}, cap), rw: new(sync.RWMutex)}
}

func (c *Channel) Insert(x interface{}) bool {
	if c.IsFull() {
		return false
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	c.ch <- x
	c.cnt += 1
	return true
}

func (c *Channel) LPop() (interface{}, bool) {
	if c.Empty() {
		return nil, false
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	var xx = <-c.ch
	c.cnt -= 1
	return xx, true
}

func (c *Channel) Empty() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnt == 0
}

func (c *Channel) IsFull() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnt >= c.cap
}
