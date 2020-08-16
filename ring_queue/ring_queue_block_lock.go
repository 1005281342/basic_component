package ring_queue

import "sync"

type RingQueueBlockRWLock struct {
	rw *sync.RWMutex
	*RingQueueBlock
}

func NewRingQueueBlockRWLock(k int) *RingQueueBlockRWLock {
	var length = k + 1
	return &RingQueueBlockRWLock{
		rw:             new(sync.RWMutex),
		RingQueueBlock: NewRingQueueBlock(length),
	}
}

func (q *RingQueueBlockRWLock) Len() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueueBlock.Len()
}

func (q *RingQueueBlockRWLock) Head() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()

	return q.RingQueueBlock.Head()
}

func (q *RingQueueBlockRWLock) Tail() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueueBlock.Tail()
}

func (q *RingQueueBlockRWLock) Insert(x interface{}) bool {
	if q.IsFull() {
		return false
	}

	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueueBlock.Insert(x)
}

func (q *RingQueueBlockRWLock) LInsert(x interface{}) bool {
	if q.IsFull() {
		return false
	}

	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueueBlock.LInsert(x)
}

func (q *RingQueueBlockRWLock) Pop() bool {
	if q.Empty() {
		return false
	}
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueueBlock.Pop()
}

func (q *RingQueueBlockRWLock) LPop() bool {
	if q.Empty() {
		return false
	}
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueueBlock.LPop()
}

// IsEmpty check the queue is empty.
func (q *RingQueueBlockRWLock) Empty() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueueBlock.Empty()
}

// IsFull check the queue is full.
func (q *RingQueueBlockRWLock) IsFull() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueueBlock.IsFull()
}
