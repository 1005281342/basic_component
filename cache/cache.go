package cache

import "time"

type BaseOp interface {
	Put(interface{}, interface{}) bool   // 添加元素到缓存中，若存在则更新元素值
	Get(interface{}) (interface{}, bool) // 从缓存中获取元素，若存在则bool == True
	Remove(interface{}) bool             // 从缓存中移除对象，若存在则返回True
}

type Cache interface {
	BaseOp
	Len() int // 当前缓存中元素个数
	Clear()   // 清空当前缓存
}

type ExpireCache interface {
	Cache
	PutWithExpire(k interface{}, v interface{}, lifeSpan time.Duration) error // 添加元素并设置存活时长
}
