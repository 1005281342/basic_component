package ring_queue

type MyCircularDeque struct {
	dq *RingQueue
}

///** Initialize your data structure here. Set the size of the deque to be k. */
//func Constructor(k int) MyCircularDeque {
//	return MyCircularDeque{dq: NewRingQueue(k)}
//}

/** Adds an item at the front of Deque. Return true if the operation is successful. */
func (this *MyCircularDeque) InsertFront(value int) bool {
	return this.dq.LInsert(value)
}

/** Adds an item at the rear of Deque. Return true if the operation is successful. */
func (this *MyCircularDeque) InsertLast(value int) bool {
	return this.dq.Insert(value)
}

/** Deletes an item from the front of Deque. Return true if the operation is successful. */
func (this *MyCircularDeque) DeleteFront() bool {
	return this.dq.LPop()
}

/** Deletes an item from the rear of Deque. Return true if the operation is successful. */
func (this *MyCircularDeque) DeleteLast() bool {
	return this.dq.Pop()
}

/** Get the front item from the deque. */
func (this *MyCircularDeque) GetFront() int {
	if this.dq.Empty() {
		return -1
	}
	return this.dq.Head().(int)
}

/** Get the last item from the deque. */
func (this *MyCircularDeque) GetRear() int {
	if this.dq.Empty() {
		return -1
	}
	return this.dq.Tail().(int)
}

/** Checks whether the circular deque is empty or not. */
func (this *MyCircularDeque) IsEmpty() bool {
	return this.dq.Empty()
}

/** Checks whether the circular deque is full or not. */
func (this *MyCircularDeque) IsFull() bool {
	return this.dq.IsFull()
}

/**
 * Your MyCircularDeque object will be instantiated and called as such:
 * obj := Constructor(k);
 * param_1 := obj.InsertFront(value);
 * param_2 := obj.InsertLast(value);
 * param_3 := obj.DeleteFront();
 * param_4 := obj.DeleteLast();
 * param_5 := obj.GetFront();
 * param_6 := obj.GetRear();
 * param_7 := obj.IsEmpty();
 * param_8 := obj.IsFull();
 */
