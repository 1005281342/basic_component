package ring_queue

type RingQueueInterface interface {
	// 获取队列中元素个数
	Len() int
	// 获取队首元素
	Head() interface{}
	// 获取队尾元素
	Tail() interface{}
	// 从队尾添加元素
	Insert(x interface{}) bool
	// 移除队首元素
	LPop() bool
	// 队列已满
	IsFull() bool
	// 队列为空
	Empty() bool
}

type RingDequeInterface interface {
	// 移除队尾元素
	Pop() bool
	// 从队首添加元素
	LInsert(x interface{}) bool
	RingQueueInterface
}
