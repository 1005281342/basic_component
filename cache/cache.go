package cache

import (
	"github.com/panjf2000/ants/v2"
	"time"
)

type BaseOp interface {
	Put(interface{}, interface{}) bool   // 添加元素到缓存中，若存在则更新元素值 返回True
	Get(interface{}) (interface{}, bool) // 从缓存中获取元素，若存在则返回True
	Remove(interface{}) bool             // 从缓存中移除对象，若存在则返回True
}

type Cache interface {
	BaseOp
	Len() int // 当前缓存中元素个数
	Clear()   // 清空当前缓存
}

type ExpireCache interface {
	Cache
	// 回收过期的元素
	DeleteExpired()
	// 不存在则添加，存在则更新
	// 需要注意永不过期与过期状态之间的切换
	PutWithExpire(k interface{}, v interface{}, lifeSpan time.Duration) bool // 添加元素并设置存活时长
}

type Opt struct {
	Callback          EvictCallback // 淘汰回调
	DefaultExpiration time.Duration // 默认过期间隔
	Interval          time.Duration // 回收间隔，限制最小为10s
	Size              int           // 缓存容量
	AntsPoolSize      int           // 协程池容量
	AntsOptionList    []ants.Option // 可选操作扩展列表
}
