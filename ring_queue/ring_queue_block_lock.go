package ring_queue

import "sync"

type RingQueueBlockLock struct {
	rw     *sync.RWMutex
	front  int           // 指向队列头部第1个有效数据的位置
	rear   int           // 指向队列尾部（即最后1个有效数据）的下一个位置，即下一个从队尾入队元素的位置
	length int           // 队列长度，非容量，slice动态库容
	nums   []interface{} // 队列元素
}

func NewRingQueueBlockLock(k int) *RingQueueBlockLock {
	var length = k + 1
	return &RingQueueBlockLock{
		rw:     new(sync.RWMutex),
		length: length,
		nums:   make([]interface{}, length),
	}
}

func (q *RingQueueBlockLock) Len() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if q.front > q.rear {
		return q.rear - q.front + q.length
	}
	return q.rear - q.front
}

func (q *RingQueueBlockLock) Head() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()

	return q.nums[q.front]
}

func (q *RingQueueBlockLock) Tail() interface{} {
	if q.Empty() {
		return nil
	}

	q.rw.RLock()
	defer q.rw.RUnlock()

	var pos int // 其实就是rear-1 //pos = (q.rear - 1 + q.length) % q.length
	if q.rear > 0 {
		pos = q.rear - 1
	} else if q.rear == 0 {
		pos = q.length - 1
	}
	return q.nums[pos]
}

func (q *RingQueueBlockLock) Insert(x interface{}) interface{} {
	if q.IsFull() {
		return false
	}

	q.rw.Lock()
	defer q.rw.Unlock()

	q.nums[q.rear] = x
	q.rear++ // q.rear = (q.rear + 1) % q.length // rear指针后移一位
	if q.rear == q.length {
		q.rear = 0
	}
	return true
}

func (q *RingQueueBlockLock) Pop() interface{} {
	if q.Empty() {
		return false
	}
	q.rw.Lock()
	defer q.rw.Unlock()
	if q.front == q.length-1 { // front指针后移一位 // q.front = (q.front + 1) % q.length
		q.front = 0
	} else {
		q.front += 1
	}
	return true
}

// IsEmpty check the queue is empty.
func (q *RingQueueBlockLock) Empty() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return q.front == q.rear
}

// IsFull check the queue is full.
func (q *RingQueueBlockLock) IsFull() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return (q.rear+1)%q.length == q.front
}
