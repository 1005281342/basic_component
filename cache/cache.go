package cache

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

type cacheType int

const (
	Simple cacheType = iota
	LRU
	LFU
	LRUk
	LRU2q
	LRUmq
	ARC
)

func NewCache(ct cacheType, opt *Opt) (ExpireCache, error) {
	switch ct {
	case Simple:
		return NewSimpleCache(opt), nil
	case LRU:
		return NewLRUCache(opt)
	case LFU:
		return NewLFUCache(opt), nil
	case LRUk:
		return NewLRUkCache(opt)
	case LRU2q:
		return NewLRU2QCache(opt)
	default:
		return nil, fmt.Errorf("not supported")
	}
}

type BaseOp interface {
	// 添加元素到缓存中，若存在则更新元素值 返回True
	Put(interface{}, interface{}) bool
	// 从缓存中获取元素，若存在则返回True
	Get(interface{}) (interface{}, bool)
	// 从缓存中移除对象，若存在则返回True
	Remove(interface{}) bool
}

type Cache interface {
	BaseOp
	// 当前缓存中元素个数
	Len() int
	// 清空当前缓存
	Clear()
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
	Callback              EvictCallback // 淘汰回调
	DefaultExpiration     time.Duration // 默认过期间隔
	Interval              time.Duration // 回收间隔，限制最小为10s
	Capacity              int           // 缓存容量
	AntsPoolCapacity      int           // 协程池容量
	AntsOptionList        []ants.Option // 可选操作扩展列表
	LruK                  int           // LRU-K/LRU-MQ的频次k
	LruKMinUpdateInterval time.Duration // LRU-K/LRU-MQ历史访问节点最小更新间隔，超过该间隔将频次置为0
	LRUMQLevel            int           //	LRUMQLevel
}
