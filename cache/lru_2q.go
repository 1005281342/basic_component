package cache

import (
	"sync"
	"time"
)

/*
1. 新访问的数据先插入到FIFO队列中；
2. 如果数据在FIFO队列中一直没有被再次访问，则最终按照FIFO规则淘汰；
3. 如果数据在FIFO队列中被再次访问，则将数据从FIFO删除，加入到LRU队列头部；
4. 如果数据在LRU队列再次被访问，则将数据移到LRU队列头部；
5. LRU队列淘汰末尾的数据。
*/
type LRU2QCache struct {
	fifo  *lru         // FIFO队列
	cache *lru         // 缓存队列
	lock  sync.RWMutex // lock
}

func NewLRU2QCache(opt *Opt) (*LRU2QCache, error) {
	var (
		fifo  *lru
		err   error
		cache *lru
	)
	if fifo, err = newLRU(opt); err != nil {
		return nil, err
	}
	if cache, err = newLRU(opt); err != nil {
		return nil, err
	}
	return &LRU2QCache{cache: cache, fifo: fifo}, nil
}

// 添加元素到缓存中，若存在则更新元素值 返回True
func (c *LRU2QCache) Put(key interface{}, value interface{}) bool {
	return c.PutWithExpire(key, value, NoExpiration)
}

// 回收过期的元素
func (c *LRU2QCache) DeleteExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 在队列中的元素并不会过期, 因此不必执行回收过期元素方法
	//c.fifo.DeleteExpired()
	c.cache.DeleteExpired()
}

// 不存在则添加，存在则更新
// 需要注意永不过期与过期状态之间的切换
func (c *LRU2QCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 1. 检查是否在FIFO队列中
	if !c.fifo.exist(key) {
		// 1.1 不存在，添加到队列中
		// 因为FIFO队列的元素值不会被查询，因此使用空结构体即可
		c.fifo.Put(key, struct{}{})
		return false
	}

	// 1.2 存在，将该元素从FIFO队列移除，然后添加元素到cache中
	c.fifo.Remove(key)
	return c.cache.PutWithExpire(key, value, lifeSpan)
}

// 从缓存中获取元素，若存在则返回True
func (c *LRU2QCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Get(key)
}

// 从缓存中移除对象，若存在则返回True
func (c *LRU2QCache) Remove(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Remove(key)
}

// 当前缓存中元素个数
func (c *LRU2QCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.Len()
}

// 清空当前缓存
func (c *LRU2QCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.fifo.Clear()
	c.cache.Clear()
}
