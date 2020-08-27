package ring_queue

import "sync"

type RingQueueRWLock struct {
	rw *sync.RWMutex
	*RingQueue
}

func NewRingQueueRWLock(k int) *RingQueueRWLock {
	var length = k + 1
	return &RingQueueRWLock{
		rw:        new(sync.RWMutex),
		RingQueue: NewRingQueue(length),
	}
}

func (q *RingQueueRWLock) Len() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueue.Len()
}

func (q *RingQueueRWLock) Head() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()

	return q.RingQueue.Head()
}

func (q *RingQueueRWLock) Tail() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueue.Tail()
}

func (q *RingQueueRWLock) Insert(x interface{}) bool {
	if q.IsFull() {
		return false
	}

	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueue.Insert(x)
}

func (q *RingQueueRWLock) LInsert(x interface{}) bool {
	if q.IsFull() {
		return false
	}

	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueue.LInsert(x)
}

func (q *RingQueueRWLock) Pop() bool {
	if q.Empty() {
		return false
	}
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueue.Pop()
}

func (q *RingQueueRWLock) LPop() bool {
	if q.Empty() {
		return false
	}
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.RingQueue.LPop()
}

// IsEmpty check the queue is empty.
func (q *RingQueueRWLock) Empty() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueue.Empty()
}

// IsFull check the queue is full.
func (q *RingQueueRWLock) IsFull() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.RingQueue.IsFull()
}
