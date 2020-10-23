package old_lfu

import (
	"sync"
	"sync/atomic"
)

type LFUCache struct {
	// 缓存存储
	cache sync.Map // map[interface{}]*lfuNode
	// 存储每个频次对应的双向链表
	freqMap sync.Map // map[uint32]*lfuNodeList
	// 缓存大小
	size uint32
	// 缓存容量
	capacity uint32
	// 当前缓存中的最小频次
	min uint32

	// lock
	lock sync.RWMutex
}

func NewLFUCache(capacity uint32) *LFUCache {
	return &LFUCache{
		cache:    sync.Map{},
		freqMap:  sync.Map{},
		size:     0,
		capacity: capacity,
		min:      0,
	}
}

// Get
func (c *LFUCache) Get(key interface{}) (interface{}, bool) {
	return c.get(key)
}
func (c *LFUCache) get(key interface{}) (interface{}, bool) {

	if node, ok := c.cache.Load(key); ok {
		var e = node.(*lfuNode)
		c.freqInc(e)
		return e.value, true
	}
	return nil, false
}

// Put
func (c *LFUCache) Put(key interface{}, value interface{}) bool {
	return c.put(key, value)
}
func (c *LFUCache) put(key interface{}, value interface{}) bool {
	if c.capacity == 0 {
		return false
	}
	var (
		newNode     *lfuNode
		oneFreqList *lfuNodeList
	)
	// 对象已存在缓存则进行更新
	if node, ok := c.cache.Load(key); ok {
		node.(*lfuNode).value = value
		c.freqInc(node.(*lfuNode))
		return false
	}

	// 若缓存容量已满, 则剔除频次最小的对象
	var evict = c.evictNode()

	newNode = newLFUNode(key, value)
	c.cache.Store(key, newNode)

	if v, ok := c.freqMap.Load(uint32(1)); !ok {
		oneFreqList = newLFUNodeList()
		c.freqMap.Store(uint32(1), oneFreqList)
	} else {
		oneFreqList = v.(*lfuNodeList)
	}

	oneFreqList.addNode(newNode)

	// c.size++
	atomic.AddUint32(&c.size, uint32(1))
	// c.min = 1
	atomic.StoreUint32(&c.min, uint32(1))
	return evict
}
func (c *LFUCache) evictNode() bool {
	if atomic.LoadUint32(&c.size) < c.capacity {
		return false
	}

	var minFreqList *lfuNodeList

	//fmt.Printf("c.size == c.capacity, c.min %d \n", c.min)
	var (
		v  interface{}
		ok bool
	)
	if v, ok = c.freqMap.Load(c.min); !ok {
		return false
	}

	// DEBUG
	//c.freqMap.Range(func(k, v interface{}) bool {
	//	fmt.Printf("iterate: %d : %#+v \n",
	//		k.(uint32),
	//		v.(*lfuNodeList).tail.prev.key)
	//	return true
	//})

	minFreqList = v.(*lfuNodeList)
	//fmt.Printf("key %#+v \n", minFreqList.tail.prev.key)
	// 从缓存中移除最小频次链表中的最后一个节点对象
	c.cache.Delete(minFreqList.tail.prev.key)
	// 从最小频次链表中移除最后一个节点
	minFreqList.removeNode(minFreqList.tail.prev)
	// c.size--
	c.sizeDec()
	return true
}

func (c *LFUCache) sizeDec() {
	atomic.AddUint32(&c.size, ^uint32(0))
}

//func (c *LFUCache) HotFixSize(d time.Duration) {
//	var timer = time.NewTicker(d)
//	select {
//	case <-timer.C:
//		c.fixSize()
//	}
//}

//func (c *LFUCache) FixSize() {
//	c.fixSize()
//}
//
//func (c *LFUCache) fixSize() {
//	var minFreqList *lfuNodeList
//	for c.size >= c.capacity {
//		//fmt.Printf("c.size == c.capacity, c.min %d \n", c.min)
//		v, _ := c.freqMap.Load(c.min)
//
//		// DEBUG
//		//c.freqMap.Range(func(k, v interface{}) bool {
//		//	fmt.Printf("iterate: %d : %#+v \n",
//		//		k.(uint32),
//		//		v.(*lfuNodeList).tail.prev.key)
//		//	return true
//		//})
//
//		minFreqList = v.(*lfuNodeList)
//		//fmt.Printf("key %#+v \n", minFreqList.tail.prev.key)
//		// 从缓存中移除最小频次链表中的最后一个节点对象
//		c.cache.Delete(minFreqList.tail.prev.key)
//		// 从最小频次链表中移除最后一个节点
//		minFreqList.removeNode(minFreqList.tail.prev)
//		// c.size--
//		atomic.AddUint32(&c.size, ^uint32(0))
//	}
//}

func (c *LFUCache) freqInc(node *lfuNode) {
	var (
		deNodeList *lfuNodeList
		ok         bool
		v          interface{}
	)

	var freq = atomic.LoadUint32(&node.freq)
	//fmt.Println("node.freq: ", node.freq)
	v, ok = c.freqMap.Load(freq)
	if !ok || v == nil {
		return
	}
	deNodeList = v.(*lfuNodeList)

	// 从双端链表中移除旧节点
	deNodeList.removeNode(node)
	if node.freq == c.min && deNodeList.head.next == deNodeList.tail {
		//atomic.StoreUint32(&c.min, freq+1)
		atomic.AddUint32(&c.min, uint32(1))
	}

	// 将新节点插入到对应频次的链表中
	node.FreqInc()
	v, ok = c.freqMap.Load(node.freq)
	if !ok {
		deNodeList = newLFUNodeList()
	} else {
		deNodeList = v.(*lfuNodeList)
	}
	c.freqMap.Store(node.freq, deNodeList)

	deNodeList.addNode(node)
}

// Remove
func (c *LFUCache) Remove(key interface{}) bool {
	var (
		v  interface{}
		ok bool
	)

	// 1. 查询元素节点
	if v, ok = c.cache.Load(key); !ok {
		// 查询的节点不存在
		return false
	}
	// 2. 根据元素节点的频率查询其所在节点链表，移除节点
	var (
		node     = v.(*lfuNode)
		freq     = atomic.LoadUint32(&node.freq)
		nodeList *lfuNodeList
	)
	if v, ok = c.freqMap.Load(freq); !ok {
		return false
	}
	nodeList = v.(*lfuNodeList)
	nodeList.removeNode(node)
	// c.size--
	c.sizeDec()
	return true
}

func (c *LFUCache) Clear() {
	// 1. 清空元素缓存Map
	c.cache.Range(func(key, _ interface{}) bool {
		c.cache.Delete(key)
		return true
	})
	// 2. 清空频次链表Map
	c.freqMap.Range(func(key, _ interface{}) bool {
		c.freqMap.Delete(key)
		return true
	})
	atomic.StoreUint32(&c.size, 0)
	atomic.StoreUint32(&c.min, 0)
}

func (c *LFUCache) GetMin() uint32 {
	return atomic.LoadUint32(&c.min)
}

func (c *LFUCache) Len() int {
	return int(atomic.LoadUint32(&c.size))
}

func (c *LFUCache) GetSize() uint32 {
	return atomic.LoadUint32(&c.size)
}

func (c *LFUCache) GetCapacity() uint32 {
	return c.capacity
}
