package old_lfu

import (
	"sync/atomic"
)

func newLFUNode(key, value interface{}) *lfuNode {
	return &lfuNode{key: key, value: value, freq: uint32(1)}
}

//func (e *lfuNode) GetValue() interface{} {
//	return e.value
//}
//
//func (e *lfuNode) UpdateValue(value interface{}) {
//	e.value = value
//}
//
//func (e *lfuNode) GetKey() interface{} {
//	return e.key
//}

func (e *lfuNode) FreqInc() {
	atomic.AddUint32(&e.freq, uint32(1))
}

func (e *lfuNode) GetFreq() uint32 {
	return atomic.LoadUint32(&e.freq)
}

type lfuNode struct {
	// Key 键
	key interface{}
	// 频次
	freq uint32

	next, prev *lfuNode

	value interface{}
}
