package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	*lru
	lock sync.RWMutex
}

func NewLRUCache(size int) (*LRUCache, error) {
	var (
		lru *lru
		err error
	)
	if lru, err = newLRU(size, nil); err != nil {
		return nil, err
	}
	return &LRUCache{lru: lru}, nil
}

func NewLRUCacheWithEvict(size int, callback EvictCallback) (*LRUCache, error) {
	var (
		lru *lru
		err error
	)
	if lru, err = newLRU(size, callback); err != nil {
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
	size      int                           // 容量
	evictList *list.List                    // 淘汰链表，需要进行淘汰时，淘汰链表尾部元素
	items     map[interface{}]*list.Element // 绑定元素key和链表节点
	onEvict   EvictCallback                 // 淘汰元素时执行的回调
}

type entry struct {
	key   interface{}
	value interface{}
}

func newLRU(size int, onEvict EvictCallback) (*lru, error) {
	if size <= 0 {
		return nil, ErrSize
	}
	c := &lru{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

// 从LRU中查找元素，返回元素值和存在标记位
func (c *lru) Get(key interface{}) (interface{}, bool) {
	var (
		node *list.Element
		ok   bool
	)
	if node, ok = c.items[key]; ok {
		c.evictList.MoveToFront(node)
		if node.Value.(*entry) == nil {
			return nil, false
		}
		return node.Value.(*entry).value, true
	}
	return nil, false
}

// Put 如果元素存在则更新, 不存在则添加；return 是否淘汰元素
func (c *lru) Put(key, value interface{}) bool {

	var (
		node *list.Element
		ok   bool
	)
	// 如果元素存在则更新
	if node, ok = c.items[key]; ok {
		c.evictList.MoveToFront(node)
		node.Value.(*entry).value = value
		return false
	}

	// 不存在则新增
	// 将元素值插入到链表头
	node = c.evictList.PushFront(&entry{key, value})
	// 绑定元素
	c.items[key] = node

	// 检查容量
	var evict = c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

// 从LRU中移除最后一个节点
func (c *lru) removeOldest() {
	var node = c.evictList.Back()
	if node != nil {
		c.removeElement(node)
	}
}

// 从LRU中移除节点；通过链表节点
func (c *lru) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}

// 从LRU中移除节点；通过key
func (c *lru) Remove(key interface{}) bool {
	var (
		node *list.Element
		ok   bool
	)
	if node, ok = c.items[key]; ok {
		c.removeElement(node)
		return true
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
}

func (c *lru) Len() int {
	return c.size
}
