package cache

import (
	"container/list"
	"runtime"
	"sync"
	"time"
)

type LFUCache struct {
	*lfu
	lock sync.RWMutex
}

func NewLFUCache(opt *Opt) *LFUCache {
	return &LFUCache{lfu: newLFU(opt)}
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

func (lc *LFUCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lfu.PutWithExpire(key, value, lifeSpan)
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
	// 过期属性
	*expire
}

func newLFU(opt *Opt) *lfu {
	var c = &lfu{
		capacity: opt.Capacity,
		onEvict:  opt.Callback,
		expire:   newExpire(opt),
		cache:    make(map[interface{}]*list.Element),
		freqMap:  make(map[int]*list.List),
	}
	if c.expire.interval > 0 {
		go c.expire.run(c)
		runtime.SetFinalizer(c.expire, stopWatchdog)
	}
	return c
}

func (c *lfu) Get(key interface{}) (interface{}, bool) {
	var (
		node *list.Element
		ok   bool
	)
	if node, ok = c.cache[key]; !ok {
		return nil, false
	}

	var et = node.Value.(*entryWithFreq)
	var value = et.item.value
	if et.Expired() {
		// 惰性回收
		// 1. 查询所在频次链表
		// 2. 从Cache中移除
		// 3. 从所在频次链表移除
		// 4. 执行回调
		c.remove(et, c.freqMap[et.freq], node)
		return nil, false
	}

	c.freqInc(node)
	return value, true
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
	return c.PutWithExpire(key, value, NoExpiration)
}

func (c *lfu) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	if c.capacity == 0 {
		return false
	}
	// 对象已存在缓存则进行更新
	if node, ok := c.cache[key]; ok {
		node.Value.(*entryWithFreq).item.value = value
		node.Value.(*entryWithFreq).item.expiration = c.absoluteTime(lifeSpan)
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
	c.cache[key] = oneFreqList.PushFront(c.newEntryWithFreq(key, value, lifeSpan))
	c.size++
	c.min = 1
	return evict
}

func (c *lfu) DeleteExpired() {
	var now = time.Now().UnixNano() // 减少系统调用
	for _, node := range c.cache {
		var et = node.Value.(*entryWithFreq)
		// 未过期
		if et.expiration <= 0 || now < et.expiration {
			continue
		}
		c.remove(et, c.freqMap[et.freq], node)
	}
}

func (c *lfu) remove(et *entryWithFreq, nodeList *list.List, node *list.Element) {
	// 移除节点
	nodeList.Remove(node)
	// 2. 从Cache中移除
	delete(c.cache, et.entry.key)
	c.size--
	// 3. 执行回调
	if c.onEvict != nil {
		_ = c.goroutinePool.Submit(func() {
			c.onEvict(et.entry.key, et.entry.item.value)
		})
	}
	entryWithFreqPool.Put(et)
}

func (c *lfu) evictNode() bool {
	if c.size < c.capacity {
		return false
	}

	var minFreqList *list.List

	var (
		v  interface{}
		ok bool
	)
	if v, ok = c.freqMap[c.min]; !ok {
		return false
	}

	minFreqList = v.(*list.List)
	// 从最小频次链表中获取最后一个节点对象
	var elem = minFreqList.Back().Value.(*entryWithFreq)
	c.remove(elem, minFreqList, minFreqList.Back())
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

	var et = node.Value.(*entryWithFreq)
	var expired = et.item.Expired()

	c.remove(et, nodeList, node)
	return !expired
}

type entryWithFreq struct {
	entry
	freq int
}

func (e *entryWithFreq) Reset() {
	e.entry.Reset()
	e.freq = 0
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

// 过期了但是未被回收也会统计在内
func (c *lfu) Len() int {
	return c.size
}

func (c *lfu) newEntryWithFreq(key, value interface{}, lifeSpan time.Duration) *entryWithFreq {
	var et = entryWithFreqPool.Get().(*entryWithFreq)
	et.Reset()
	et.freq = 1
	et.entry.key = key
	et.entry.item.value = value
	et.entry.item.expiration = c.absoluteTime(lifeSpan)
	return et
}
