package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
新插入的数据放入Q0；
每个队列按照LRU管理数据，再次访问的数据移动到头部；
当数据的访问次数达到一定次数，需要提升优先级时，将数据从当前队列删除，加入到高一级队列的头部；
为了防止高优先级数据永远不被淘汰，当数据在指定的时间里没有被访问时，需要降低优先级，将数据从当前队列删除，加入到低一级的队列头部；
需要淘汰数据时，从最低一级队列开始按照LRU淘汰；每个队列淘汰数据时，将数据从缓存中删除，将数据索引加入Q-history头部；
如果数据在Q-history中被重新访问，则重新计算其优先级，移到目标队列的头部；
Q-history按照LRU淘汰数据的索引。
*/

type LRUMQCache struct {
	cache     *lru         // MQ Cache
	levelList []*lru       // 将历史访问队列记为levelList[0], level list
	level     int          // 级
	capacity  int          // 容量
	size      int          // 使用大小
	lock      sync.RWMutex // lock
	k         int          // 升级
	ik        int64        // 服务绝对自增k值
}

// cache如何与levelList关联？？ 比如从Cache中移除一个key，如何保证正确地从对应的level lru移除节点

func NewLRUMQCache(opt *Opt) (*LRUMQCache, error) {
	var (
		cache *lru
		err   error
	)
	if cache, err = newLRU(opt); err != nil {
		return nil, err
	}
	var levelList = make([]*lru, opt.LRUMQLevel+1, opt.LRUMQLevel+1)
	for i := 0; i < opt.Capacity; i++ {
		if levelList[i], err = newLRU(opt); err != nil {
			return nil, err
		}
	}
	var c = &LRUMQCache{cache: cache, levelList: levelList, level: opt.LRUMQLevel, k: opt.LruK, capacity: opt.Capacity}
	// 绝对K值自增。
	// 在一个自增间隔添加的元素在一个优先级队列中
	go func(d time.Duration, c *LRUMQCache) {
		var tc = time.NewTicker(d)
		for range tc.C {
			atomic.AddInt64(&c.ik, 1)
		}
	}(opt.LruKMinUpdateInterval, c)
	// 定期更新等级。100ms
	go c.updateLevel()
	return c, nil
}

func (c *LRUMQCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var v, ok = c.cache.Get(key)
	if !ok {
		return nil, false
	}

	var mq = v.(*mqEntry)
	mq.freq++
	// 惰性检查级别
	mq.level = c.checkLevel(key, mq.freq, mq.level)
	return mq.value, true
}

func (c *LRUMQCache) updateLevel() {
	var tc = time.NewTicker(100 * time.Millisecond)
	for range tc.C {
		for i := range c.levelList {
			c.lock.Lock()
			for node := c.levelList[i].evictList.Back(); node != nil; node = node.Prev() {
				var mqe = node.Value.(*mqEntry)
				mqe.level = c.checkLevel(mqe.key, mqe.freq, mqe.level)
			}
			// 如果在中间过程出现移除退出，则可能出现死锁
			c.lock.Unlock()
		}
	}
}

// return [0, level], [0, level+1)
func (c *LRUMQCache) checkLevel(key interface{}, freq int, level int) int {
	var a = atomic.LoadInt64(&c.ik)
	var newLevel = (freq - int(a)) / c.k
	if level != newLevel {
		// 从旧级别中移除节点
		c.levelList[level].Remove(key)
		// 添加到新级别中
		c.levelList[newLevel].putItem(key, struct{}{})
	}
	return newLevel
}

func (c *LRUMQCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 1. 查询是否存在，存在则更新
	// 	1.1 检查level
	// 2. 不存在则新增，其level为1
	//  2.1 是否有数据被逐出

	var (
		_, ok = c.cache.Get(key)
		elem  interface{}
		evict bool
	)
	if !ok {
		var et = mqEntryPool.Get().(*mqEntry)
		et.Reset()
		et.key = key
		et.level = 1
		et.freq = int(c.ik)
		et.value = value
		et.expiration = c.cache.absoluteTime(lifeSpan)
		elem, evict = c.cache.putItem2(key, et)

		// 如果淘汰了元素则将等级队列中对应的元素移除
		if evict && elem != nil {
			var el = elem.(*mqEntry)
			c.levelList[el.freq].Remove(key)
		}
		return evict
	}

	var mqe = c.cache.items[key].Value.(*mqEntry)
	// 移除等级队列中的元素
	c.levelList[mqe.level].Remove(key)
	// 更新
	mqe.key = key
	mqe.value = value
	mqe.freq = int(c.ik)
	mqe.level = 1
	// 注册到等级队列中
	if elem, evict = c.levelList[mqe.level].putItem2(key, struct{}{}); evict && elem != nil {
		//如果淘汰了数据
		c.cache.Remove(key)
	}
	return false
}

func (c *LRUMQCache) Put(key, value interface{}) bool {
	return c.PutWithExpire(key, value, NoExpiration)
}

func (c *LRUMQCache) Remove(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	var elem, ok = c.cache.items[key]
	if !ok {
		return false
	}

	var mqe = elem.Value.(*mqEntry)
	// 1. 从Cache中移除
	ok = c.cache.Remove(key)
	// 2. 从等级队列中移除
	c.levelList[mqe.level].Remove(key)
	return ok
}

func (c *LRUMQCache) DeleteExpired() {

}

func (c *LRUMQCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.Len()
}

func (c *LRUMQCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Clear()
	for i := range c.levelList {
		c.levelList[i].Clear()
	}
}

type mqEntry struct {
	entry
	freq  int
	level int
}

func (m *mqEntry) Reset() {
	m.entry.Reset()
	m.freq = 0
	m.level = 0
}
