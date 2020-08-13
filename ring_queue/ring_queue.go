package ring_queue

type RingQueue struct {
	cap   int
	queue []interface{}
	index int // tail 队尾
	head  int // head 队首
}

func NewRingQueue(cap int) *RingQueue {

	if cap <= 1 {
		return nil
	}

	return &RingQueue{
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
func (r *RingQueue) Len() int {
	if r.head > r.index {
		return r.index - r.head + 1 + r.cap
	}
	return r.index - r.head + 1
}

// 从队尾插入一个元素，如果队列满了则弹出队首元素
func (r *RingQueue) Insert(x interface{}) interface{} {
	var node = r.Head()
	var full = r.IsFull()

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
func (r *RingQueue) LInsert(x interface{}) interface{} {

	// 如果没有满，则直接在队首插入元素即可
	if !r.IsFull() {
		if r.head == 0 {
			r.head = r.cap
		}
		r.head -= 1
		r.queue[r.head] = x
		// 队列未满，无须弹出元素
		return nil
	}

	var node = r.Tail()

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
func (r *RingQueue) Head() interface{} {
	return r.queue[r.head]
}

// 获取队尾元素
func (r *RingQueue) Tail() interface{} {
	return r.queue[r.index]
}

// 弹出队首元素
func (r *RingQueue) LPop() interface{} {
	if r.Empty() {
		return nil
	}

	var node = r.Head()
	if r.head == r.cap-1 {
		r.head = 0
	} else {
		r.head += 1
	}
	return node
}

// 弹出队尾元素
func (r *RingQueue) Pop() interface{} {
	if r.Empty() {
		return nil
	}

	var node = r.Tail()
	if r.index == 0 {
		r.index = r.cap - 1
	} else {
		r.index -= 1
	}
	return node
}

// 队列已经满了
func (r *RingQueue) IsFull() bool {
	return (r.index+1)%r.cap == r.head
}

func (r *RingQueue) Empty() bool {
	return r.head == r.index
}
