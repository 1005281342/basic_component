package cache

import (
	"container/list"
	"sync"
)

type LFUCache struct {
	*lfu
	lock sync.RWMutex
}

func NewLFUCache(size int) *LFUCache {
	return &LFUCache{lfu: newLFU(size, nil)}
}
func NewLFUCacheWithCallBack(size int, callback EvictCallback) *LFUCache {
	return &LFUCache{lfu: newLFU(size, callback)}
}

func (lc *LFUCache) Get(key interface{}) (interface{}, bool) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lfu.Get(key)
}

func (lc *LFUCache) Put(key, value interface{}) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lfu.Put(key, value)
}

func (lc *LFUCache) Remove(key interface{}) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lfu.Remove(key)
}

func (lc *LFUCache) Len() int {
	lc.lock.RLock()
	defer lc.lock.RUnlock()
	return lc.lfu.Len()
}

func (lc *LFUCache) Clear() {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	lc.lfu.Clear()
}

type lfu struct {
	// 缓存存储
	cache map[interface{}]*list.Element
	// 存储每个频次对应的双向链表
	freqMap map[int]*list.List
	// 缓存大小
	size int
	// 缓存容量
	capacity int
	// 当前缓存中的最小频次
	min int
	// 淘汰元素时执行的回调
	onEvict EvictCallback
}

func newLFU(size int, callback EvictCallback) *lfu {
	return &lfu{size: size, onEvict: callback}
}

func (c *lfu) Get(key interface{}) (interface{}, bool) {
	if node, ok := c.cache[key]; ok {
		c.freqInc(node)
		return node.Value.(*list.Element).Value, true
	}
	return nil, false
}

func (c *lfu) freqInc(node *list.Element) {
	var (
		deNodeList *list.List
		ok         bool
		v          interface{}
	)

	var freq = node.Value.(*entryWithFreq).freq
	//fmt.Println("node.freq: ", node.freq)
	v, ok = c.freqMap[freq]
	if !ok || v == nil {
		return
	}
	deNodeList = v.(*list.List)

	// 从双端链表中移除旧节点
	deNodeList.Remove(node)
	if freq == c.min && deNodeList.Len() == 0 {
		c.min = 1
	}

	// 将新节点插入到对应频次的链表中
	freq++
	// 更新频率
	node.Value.(*entryWithFreq).freq = freq
	if deNodeList, ok = c.freqMap[freq]; !ok {
		deNodeList = list.New()
		c.freqMap[freq] = deNodeList
	}
	deNodeList.PushFront(node)
}

func (c *lfu) Put(key interface{}, value interface{}) bool {
	if c.capacity == 0 {
		return false
	}
	// 对象已存在缓存则进行更新
	if node, ok := c.cache[key]; ok {
		node.Value.(*entryWithFreq).value = value
		node.Value.(*entryWithFreq).freq++
		return false
	}

	var (
		// 若缓存容量已满, 则剔除频次最小的对象
		evict       = c.evictNode()
		oneFreqList *list.List
		ok          bool
	)

	if oneFreqList, ok = c.freqMap[1]; !ok {
		oneFreqList = list.New()
		c.freqMap[1] = oneFreqList
	}
	c.cache[key] = oneFreqList.PushFront(newEntryWithFreq(key, value))
	c.size++
	c.min = 1
	return evict
}

func (c *lfu) evictNode() bool {
	if c.size < c.capacity {
		return false
	}

	var minFreqList *list.List

	//fmt.Printf("c.size == c.capacity, c.min %d \n", c.min)
	var (
		v  interface{}
		ok bool
	)
	if v, ok = c.freqMap[c.min]; !ok {
		return false
	}

	// DEBUG
	//c.freqMap.Range(func(k, v interface{}) bool {
	//	fmt.Printf("iterate: %d : %#+v \n",
	//		k.(uint32),
	//		v.(*lfuNodeList).tail.prev.key)
	//	return true
	//})

	minFreqList = v.(*list.List)
	//fmt.Printf("key %#+v \n", minFreqList.tail.prev.key)
	// 从缓存中移除最小频次链表中的最后一个节点对象
	//c.cache.Delete(minFreqList.tail.prev.key)
	var elem = minFreqList.Back().Value.(*entryWithFreq)
	if c.onEvict != nil {
		go c.onEvict(elem.key, elem.value)
	}
	delete(c.cache, elem.key)
	// 从最小频次链表中移除最后一个节点
	//minFreqList.removeNode(minFreqList.tail.prev)
	minFreqList.Remove(minFreqList.Back())
	c.size--
	//c.sizeDec()
	return true
}

// Remove
func (c *lfu) Remove(key interface{}) bool {
	var (
		node *list.Element
		ok   bool
	)

	// 1. 查询元素节点
	if node, ok = c.cache[key]; !ok {
		// 查询的节点不存在
		return false
	}
	// 2. 根据元素节点的频率查询其所在节点链表，移除节点
	var (
		freq     = node.Value.(*entryWithFreq).freq
		nodeList *list.List
	)
	if nodeList, ok = c.freqMap[freq]; !ok {
		return false
	}
	nodeList.Remove(node)
	c.size--
	return true
}

type entryWithFreq struct {
	*entry
	freq int
}

func (c *lfu) Clear() {
	// 1. 清空元素缓存Map
	for k, v := range c.cache {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entryWithFreq).value)
		}
		delete(c.cache, k)
	}
	// 2. 清空频次链表Map
	for k, v := range c.freqMap {
		v.Init()
		delete(c.freqMap, k)
	}
	c.size = 0
	c.min = 0
}

func (c *lfu) Len() int {
	return c.size
}

func newEntryWithFreq(key, value interface{}) *entryWithFreq {
	return &entryWithFreq{&entry{key: key, value: value}, 1}
}
