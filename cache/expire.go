package cache

import (
	"time"
)

const (
	// 永不过期
	NoExpiration time.Duration = -1
	// 最小过期时间阈值
	DefaultExpirationThreshold time.Duration = 0
)

type expire struct {
	defaultExpiration time.Duration // 默认多长时间过期
	watchdog                        // 看门狗，定期回收过期元素

	// 协程池
	goroutinePool
}

// 获取绝对时间
func (e *expire) absoluteTime(d time.Duration) int64 {
	var t int64
	if d > 0 {
		t = time.Now().Add(d).UnixNano()
	}
	return t
}

// defaultExpiration 默认过期时间，interval看门狗回收间隔
func newExpire(opt *Opt) *expire {
	// 回收间隔不规范
	if opt.Interval <= 0 {
		// 参考redis设置定期回收间隔为10s
		opt.Interval = time.Second * 10
	}
	return &expire{
		defaultExpiration: opt.DefaultExpiration,
		watchdog:          watchdog{stop: make(chan struct{}), interval: opt.Interval},
		goroutinePool:     newGoroutinePool(opt.AntsPoolCapacity, opt.AntsOptionList...),
	}
}
