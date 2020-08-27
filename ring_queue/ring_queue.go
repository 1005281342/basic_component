package ring_queue

type RingQueue struct {
	front  int           // 指向队列头部第1个有效数据的位置
	rear   int           // 指向队列尾部（即最后1个有效数据）的下一个位置，即下一个从队尾入队元素的位置
	length int           // 队列长度，非容量，slice动态库容
	nums   []interface{} // 队列元素
}

// NewRingQueue Initialize the structure. Set the size of the queue to be k.
func NewRingQueue(k int) *RingQueue {
	var length = k + 1
	return &RingQueue{
		length: length,
		nums:   make([]interface{}, length),
	}
}

func (q *RingQueue) Len() int {
	if q.front > q.rear {
		return q.rear - q.front + q.length
	}
	return q.rear - q.front
}

func (q *RingQueue) Head() interface{} {
	if q.Empty() {
		return nil
	}

	return q.nums[q.front]
}

func (q *RingQueue) Tail() interface{} {
	if q.Empty() {
		return nil
	}
	var pos int // 其实就是rear-1 //pos = (q.rear - 1 + q.length) % q.length
	if q.rear > 0 {
		pos = q.rear - 1
	} else if q.rear == 0 {
		pos = q.length - 1
	}
	return q.nums[pos]
}

func (q *RingQueue) LInsert(x interface{}) bool {
	if q.IsFull() {
		return false
	}
	if q.front == 0 {
		q.front = q.length - 1
	} else {
		q.front -= 1
	}
	q.nums[q.front] = x
	return true
}

func (q *RingQueue) Insert(x interface{}) bool {
	if q.IsFull() {
		return false
	}
	q.nums[q.rear] = x
	q.rear++ // q.rear = (q.rear + 1) % q.length // rear指针后移一位
	if q.rear == q.length {
		q.rear = 0
	}
	return true
}

func (q *RingQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	if q.rear == 0 {
		q.rear = q.length - 1
	} else {
		q.rear -= 1
	}
	return true
}

func (q *RingQueue) LPop() bool {
	if q.Empty() {
		return false
	}

	if q.front == q.length-1 { // front指针后移一位 // q.front = (q.front + 1) % q.length
		q.front = 0
	} else {
		q.front += 1
	}
	return true
}

// IsEmpty check the queue is empty.
func (q *RingQueue) Empty() bool {
	return q.front == q.rear
}

// IsFull check the queue is full.
func (q *RingQueue) IsFull() bool {
	return (q.rear+1)%q.length == q.front
}
