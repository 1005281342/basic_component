package ring_queue

type MyCircularQueue struct {
	rq *RingQueueBlock
}

/** Initialize your data structure here. Set the size of the queue to be k. */
func Constructor(k int) MyCircularQueue {
	return MyCircularQueue{rq: NewRingQueueBlock(k)}
}

/** Insert an element into the circular queue. Return true if the operation is successful. */
func (this *MyCircularQueue) EnQueue(value int) bool {
	var b = this.rq.IsFull()
	this.rq.Insert(value)
	return !b
}

/** Delete an element from the circular queue. Return true if the operation is successful. */
func (this *MyCircularQueue) DeQueue() bool {
	if this.rq.Empty() {
		return false
	}
	this.rq.Pop()
	return true
}

/** Get the front item from the queue. */
func (this *MyCircularQueue) Front() int {
	if this.rq.Empty() {
		return -1
	}
	return this.rq.Head().(int)
}

/** Get the last item from the queue. */
func (this *MyCircularQueue) Rear() int {
	if this.rq.Empty() {
		return -1
	}
	return this.rq.Tail().(int)
}

/** Checks whether the circular queue is empty or not. */
func (this *MyCircularQueue) IsEmpty() bool {
	return this.rq.Empty()
}

/** Checks whether the circular queue is full or not. */
func (this *MyCircularQueue) IsFull() bool {
	return this.rq.IsFull()
}

/**
 * Your MyCircularQueue object will be instantiated and called as such:
 * obj := Constructor(k);
 * param_1 := obj.EnQueue(value);
 * param_2 := obj.DeQueue();
 * param_3 := obj.Front();
 * param_4 := obj.Rear();
 * param_5 := obj.IsEmpty();
 * param_6 := obj.IsFull();
 */
