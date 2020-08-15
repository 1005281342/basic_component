package ring_queue

import "sync"

type RingQueueRWLock struct {
	rw    *sync.RWMutex
	cap   int
	queue []interface{}
	index int // tail 队尾
	head  int // head 队首
}

func NewRingQueueRWLock(cap int) *RingQueueRWLock {

	if cap <= 1 {
		return nil
	}

	return &RingQueueRWLock{
		rw:    new(sync.RWMutex),
		cap:   cap,
		queue: make([]interface{}, cap, cap),
		index: 0,
		head:  0,
	}
}

// 获取当前队列元素个数
// 1, 2, 3, 4, 5, 6
// tail 5, head 0, tail - head + 1 = 6
// tail 1, head 3, tail - head + 1 + cap = 5
func (r *RingQueueRWLock) Len() int {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.head > r.index {
		return r.index - r.head + 1 + r.cap
	}
	return r.index - r.head + 1
}

// 从队尾插入一个元素，如果队列满了则弹出队首元素
func (r *RingQueueRWLock) Insert(x interface{}) interface{} {

	var node = r.Head()
	var full = r.IsFull()

	r.rw.Lock()
	defer r.rw.Unlock()

	r.index++
	if r.index == r.cap {
		r.index = 0
	}

	r.queue[r.index] = x

	// 维护队首
	if full {
		r.head = (r.index + 1) % r.cap
	}
	return node
}

// 从队首插入一个元素，如果队列满了则弹出队尾元素
// 该方法适用于：如临时有一个任务需要优先执行
func (r *RingQueueRWLock) LInsert(x interface{}) interface{} {

	// 如果没有满，则直接在队首插入元素即可
	if !r.IsFull() {

		r.rw.Lock()
		defer r.rw.Unlock()

		if r.head == 0 {
			r.head = r.cap
		}
		r.head -= 1
		r.queue[r.head] = x
		// 队列未满，无须弹出元素
		return nil
	}

	var node = r.Tail()
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.index == 0 {
		r.head = 0
		r.queue[0] = x
		r.index = r.cap - 1
		return node
	}

	// 处理队首
	r.head = r.index
	r.queue[r.head] = x

	// 处理队尾
	r.index -= 1
	return node
}

// 获取队首元素
func (r *RingQueueRWLock) Head() interface{} {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.queue[r.head]
}

// 获取队尾元素
func (r *RingQueueRWLock) Tail() interface{} {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.queue[r.index]
}

// 弹出队首元素
func (r *RingQueueRWLock) LPop() interface{} {

	if r.Empty() {
		return nil
	}

	var node = r.Head()

	r.rw.Lock()
	defer r.rw.Unlock()

	if r.head == r.cap-1 {
		r.head = 0
	} else {
		r.head += 1
	}
	return node
}

// 弹出队尾元素
func (r *RingQueueRWLock) Pop() interface{} {
	if r.Empty() {
		return nil
	}

	var node = r.Tail()
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.index == 0 {
		r.index = r.cap - 1
	} else {
		r.index -= 1
	}
	return node
}

// 队列已经满了
func (r *RingQueueRWLock) IsFull() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return (r.index+1)%r.cap == r.head
}

func (r *RingQueueRWLock) Empty() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.head == r.index
}
