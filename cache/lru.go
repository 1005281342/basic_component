package cache

import (
	"container/list"
	"runtime"
	"sync"
	"time"
)

type LRUCache struct {
	*lru
	lock sync.RWMutex
}

func NewLRUCache(opt *Opt) (*LRUCache, error) {
	var (
		lru *lru
		err error
	)
	if lru, err = newLRU(opt); err != nil {
		return nil, err
	}
	return &LRUCache{lru: lru}, nil
}

func (lc *LRUCache) Get(key interface{}) (interface{}, bool) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lru.Get(key)
}

func (lc *LRUCache) Put(key, value interface{}) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lru.Put(key, value)
}

func (lc *LRUCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lru.PutWithExpire(key, value, lifeSpan)
}

func (lc *LRUCache) Remove(key interface{}) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.lru.Remove(key)
}

func (lc *LRUCache) Len() int {
	lc.lock.RLock()
	defer lc.lock.RUnlock()
	return lc.lru.Len()
}

func (lc *LRUCache) Clear() {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	lc.lru.Clear()
}

type lru struct {
	capacity  int                           // 缓存容量
	size      int                           // 使用节点
	evictList *list.List                    // 淘汰链表，需要进行淘汰时，淘汰链表尾部元素
	items     map[interface{}]*list.Element // 绑定元素key和链表节点
	onEvict   EvictCallback                 // 淘汰元素时执行的回调
	*expire                                 // 过期属性
}

type entry struct {
	key interface{}
	item
}

func (e *entry) Reset() {
	e.item.Reset()
	e.key = nil
}

func newLRU(opt *Opt) (*lru, error) {
	if opt.Capacity <= 0 {
		return nil, ErrSize
	}
	c := &lru{
		capacity:  opt.Capacity,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   opt.Callback,
		expire:    newExpire(opt),
	}
	if c.expire.interval > 0 {
		go c.expire.run(c)
		runtime.SetFinalizer(c.expire, stopWatchdog)
	}
	return c, nil
}

// 从LRU中查找元素，返回元素值和存在标记位
func (c *lru) Get(key interface{}) (interface{}, bool) {
	var (
		node *list.Element
		ok   bool
	)
	if node, ok = c.items[key]; !ok {
		return nil, false
	}

	var et = node.Value.(*entry)
	if et == nil {
		return nil, false
	}

	if et.Expired() {
		c.removeElement(node)
		return nil, false
	}

	c.evictList.MoveToFront(node)
	return et.item.value, true
}

func (c *lru) exist(key interface{}) bool {
	var _, ok = c.items[key]
	return ok
}

// Put 如果元素存在则更新, 不存在则添加；return 是否淘汰元素
func (c *lru) Put(key, value interface{}) bool {
	return c.PutWithExpire(key, value, NoExpiration)
}

func (c *lru) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {

	var (
		node *list.Element
		ok   bool
	)
	// 如果元素存在则更新
	if node, ok = c.items[key]; ok {
		c.evictList.MoveToFront(node)
		node.Value.(*entry).item.value = value
		node.Value.(*entry).item.expiration = c.absoluteTime(lifeSpan)
		return false
	}

	return c.put(key, value, lifeSpan)
}

func (c *lru) put(key, value interface{}, lifeSpan time.Duration) bool {
	// 不存在则新增
	// 将元素值插入到链表头
	var et = entryPool.Get().(*entry)
	et.Reset()
	et.key = key
	et.item.value = value
	et.item.expiration = c.absoluteTime(lifeSpan)
	return c.putItem(key, et)
}

func (c *lru) putItem(key, it interface{}) bool {
	var node = c.evictList.PushFront(it)
	// 绑定元素
	c.items[key] = node
	c.size++

	// 检查容量
	var evict = c.evictList.Len() > c.capacity
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *lru) putItem2(key, it interface{}) (interface{}, bool) {
	var node = c.evictList.PushFront(it)
	// 绑定元素
	c.items[key] = node
	c.size++

	// 检查容量
	var evict = c.evictList.Len() > c.capacity
	if evict {
		return c.removeOldest(), true
	}
	return nil, false
}

func (c *lru) DeleteExpired() {
	var now = time.Now().UnixNano() // 减少系统调用
	for node := c.evictList.Front(); node != nil; node = node.Next() {
		var it = node.Value.(*entry).item
		if it.expiration > 0 && now > it.expiration {
			c.removeElement(node)
		}
	}
}

// 从LRU中移除最后一个节点
func (c *lru) removeOldest() interface{} {
	var node = c.evictList.Back()
	if node != nil {
		return c.removeElement(node)
	}
	return nil
}

// 从LRU中移除节点；通过链表节点
func (c *lru) removeElement(e *list.Element) interface{} {
	var elem = c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	c.size--
	if c.onEvict != nil {
		_ = c.goroutinePool.Submit(func() {
			c.onEvict(kv.key, kv.item.value)
		})
	}
	entryPool.Put(kv)
	return elem
}

// 从LRU中移除节点；通过key
func (c *lru) Remove(key interface{}) bool {
	var (
		node *list.Element
		ok   bool
	)
	if node, ok = c.items[key]; ok {

		var expired = node.Value.(*entry).Expired()
		c.removeElement(node)
		return !expired
	}
	return false
}

func (c *lru) Clear() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
	c.size = 0
}

func (c *lru) Len() int {
	return c.size
}
