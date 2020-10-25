package cache

import (
	"runtime"
	"sync"
	"time"
)

type SimpleCache struct {
	*simple
	lock sync.RWMutex
}

func NewSimpleCache(opt *Opt) *SimpleCache {
	return &SimpleCache{simple: newSimple(opt)}
}

func (sc *SimpleCache) Put(key interface{}, value interface{}) bool {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	return sc.simple.Put(key, value)
}

func (sc *SimpleCache) PutWithExpire(key interface{}, value interface{}, lifeSpan time.Duration) bool {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	return sc.simple.PutWithExpire(key, value, lifeSpan)
}

func (sc *SimpleCache) Get(key interface{}) (interface{}, bool) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	return sc.simple.Get(key)
}

func (sc *SimpleCache) Remove(key interface{}) bool {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	return sc.simple.Remove(key)
}

func (sc *SimpleCache) Len() int {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	return sc.simple.Len()
}

func (sc *SimpleCache) Clear() {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	sc.simple.Clear()
}

type simple struct {
	size    int
	items   map[interface{}]*item
	onEvict EvictCallback // 淘汰元素时执行的回调
	*expire               // 过期属性
}

func newSimple(opt *Opt) *simple {
	var s = &simple{
		items:   make(map[interface{}]*item),
		onEvict: opt.Callback,
		expire:  newExpire(opt),
	}
	if s.expire.interval > 0 {
		go s.expire.run(s)
		runtime.SetFinalizer(s.expire, stopWatchdog)
	}
	return s
}

// Put 添加元素到缓存中，若存在则更新元素值
func (s *simple) Put(key interface{}, value interface{}) bool {
	return s.PutWithExpire(key, value, NoExpiration)
}
func (s *simple) PutWithExpire(k interface{}, v interface{}, lifeSpan time.Duration) bool {
	var (
		it *item
		ok bool
	)

	// 过期阈值校验
	if lifeSpan == DefaultExpirationThreshold {
		lifeSpan = s.defaultExpiration
	}

	// 存在于缓存中
	if it, ok = s.items[k]; ok {
		var add = it.Expired()
		it.value = v
		it.expiration = s.absoluteTime(lifeSpan)
		return add
	}

	// 不存在，新增
	it = itemPool.Get().(*item)
	// 清空
	it.Reset()
	// 赋值
	it.value = v
	it.expiration = s.absoluteTime(lifeSpan)
	s.items[k] = it
	s.size++
	return true
}

// 从缓存中获取元素，若存在则bool == True
func (s *simple) Get(key interface{}) (interface{}, bool) {

	var (
		it *item
		ok bool
	)

	// 先查询是否存在
	if it, ok = s.items[key]; !ok {
		return nil, false
	}

	// 判断是否过期
	// 过期则触发惰性回收
	if it.Expired() {
		// 惰性回收
		s.remove(key, it)
		return nil, false
	}

	// 返回值
	return it.value, true
}

// 从缓存中移除对象，若存在则返回True
func (s *simple) Remove(key interface{}) bool {

	var (
		it      *item
		ok      bool
		expired bool
	)

	// 先查询是否存在
	if it, ok = s.items[key]; !ok {
		return false
	}
	expired = it.Expired()
	if !expired {
		s.size--
	}
	s.remove(key, it)
	return !expired
}

func (s *simple) remove(key interface{}, it *item) {
	// 惰性回收
	var val = it.value
	delete(s.items, key)
	_ = s.goroutinePool.Submit(func() {
		callEvict(s.onEvict, key, val)
	})
	itemPool.Put(it)
}

func (s *simple) Len() int {
	return s.size
}

func (s *simple) Clear() {
	for k, v := range s.items {
		s.remove(k, v)
	}
	s.size = 0
}

// 回收过期的元素
func (s *simple) DeleteExpired() {
	var now = time.Now().UnixNano() // 减少系统调用
	for k, v := range s.items {
		if v.expiration > 0 && now > v.expiration {
			s.remove(k, v)
		}
	}
}

type item struct {
	value      interface{} // 元素值
	expiration int64       // 绝对过期时间
}

func (i *item) Reset() {
	i.value = nil
	i.expiration = 0
}

// Expired 是否过期
func (i item) Expired() bool {
	if i.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.expiration
}

// Value 查询元素值，若不存在返回 nil, false
func (i item) Value() (interface{}, bool) {
	if i.Expired() {
		return nil, false
	}
	return i.value, true
}
