package cache

import (
	"time"
)

type watchdog struct {
	interval time.Duration
	stop     chan struct{}
}

// 启动看门狗
func (w *watchdog) run(c ExpireCache) {
	var ticker = time.NewTicker(w.interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-w.stop:
			ticker.Stop()
			return
		}
	}
}

// 查询定期回收间隔
func (w *watchdog) Interval() time.Duration {
	return w.interval
}

func stopWatchdog(c *expire) {
	c.stop <- struct{}{}
}
