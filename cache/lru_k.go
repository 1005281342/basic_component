package cache

import (
	"container/list"
	"sync"
	"time"
)

/*
1. 数据第一次被访问，加入到访问历史记录表（简称记录表）；在记录表中对应的K单元中设置最后访问时间=new()，且设置访问次数为1；
2. 如果数据访问次数没有达到K次，则访问次数+1。最后访问时间与当前时间间隔超过预设的值(如30秒)则该热点值做衰减，再累加1；
3. 当数据访问计数超过(>=)K次后，则访问次数+1。将数据保存到LRU缓存队列中，缓存队列重新按照时间排序；
4. LRU缓存队列中数据被再次访问后，重新排序；
5. LRU缓存队列需要淘汰数据时，淘汰缓存队列中排在末尾的数据，即：淘汰“倒数第K次访问离现在最久”的数据。
*/

type LRUkCache struct {
	// 历史访问队列
	history *lru
	// 缓存
	cache *lru
	// lock
	lock sync.RWMutex
	// 写缓存频次
	k int
}

func NewLRUkCache(opt *Opt) (*LRUkCache, error) {
	if opt.LruKMinUpdateInterval == 0 {
		opt.LruKMinUpdateInterval = DefaultLruKMinUpdateInterval
	}

	var (
		history *lru
		err     error
	)
	if history, err = newLRU(opt); err != nil {
		return nil, err
	}
	var cache, _ = newLRU(opt)
	return &LRUkCache{k: opt.LruK, history: history, cache: cache}, nil
}

func (c *LRUkCache) Put(key interface{}, value interface{}) bool {
	return c.PutWithExpire(key, value, NoExpiration)
}

// 添加元素到缓存中，若存在则更新元素值 返回True
func (c *LRUkCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {

	var (
		it *list.Element
		ok bool
	)

	c.lock.Lock()
	defer c.lock.Unlock()
	// 1. 是否已存在于缓存中
	if it, ok = c.cache.items[key]; ok {
		it.Value.(*entry).item.value = value
		c.cache.evictList.MoveToFront(it)
		return false
	}

	var het = entryWithHistoryPool.Get().(*entryWithHistory)
	het.Reset()
	// 2. 检查是否在历史访问队列中
	// 2.1 节点存在
	if it, ok = c.history.items[key]; ok {
		het = it.Value.(*entryWithHistory)
		// 热度削减
		het.Hot()

		// 访问频次自增
		het.freq++
		it.Value = het

		// 频次达到条件
		if het.freq >= c.k {
			// 从历史访问列表中移除
			c.history.removeElement(it)

			// 添加到缓存中
			return c.cache.put(key, value, lifeSpan)
		}
		// 频次未达条件, 调整历史访问列表
		c.history.evictList.MoveToFront(it)
		return false
	}

	// 2.2 不存在于历史访问列表中
	// 记录key
	het.key = key
	//het.entry.key = key
	//het.entry.item.value = value
	het.freq = 1
	// 更新时间
	het.updateTime = time.Now().UnixNano()
	c.history.putItem(key, het)
	return false
}

// 从缓存中获取元素，若存在则返回True
func (c *LRUkCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Get(key)
}

// 从缓存中移除对象，若存在则返回True
func (c *LRUkCache) Remove(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Remove(key)
}

// 回收过期的元素
func (c *LRUkCache) DeleteExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.DeleteExpired()
	c.history.DeleteExpired()
}

func (c *LRUkCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.Len()
}

func (c *LRUkCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Clear()
	c.history.Clear()
}

type entryWithHistory struct {
	key        interface{}
	freq       int   // 频次
	updateTime int64 // 更新绝对时间
}

func (e *entryWithHistory) Reset() {
	e.key = nil
	e.freq = 0
	e.updateTime = 0
}

func (e *entryWithHistory) Hot() {
	var now = time.Now().UnixNano()
	var cnt = float64(now-e.updateTime) / float64(DefaultLruKMinUpdateInterval)
	e.freq -= int(cnt)
	if e.freq < 0 {
		e.freq = 0
	}
}

const (
	DefaultLruKMinUpdateInterval = 10 * time.Second
)
