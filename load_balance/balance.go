package load_balance

import "sync"

// balance 负载均衡基础结构体
type balance struct {
	// RWLock
	rwLock sync.RWMutex

	// 后面统一集成服务发现功能
}
