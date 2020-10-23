package old_lfu

import (
	"sync"
)

type lfuNodeList struct {
	// 头哨兵
	head *lfuNode
	// 尾哨兵
	tail *lfuNode
	// 读写锁
	lock sync.RWMutex
}

func newLFUNodeList() *lfuNodeList {

	var (
		head *lfuNode
		tail *lfuNode
	)

	head = newLFUNode(struct{}{}, struct{}{})
	tail = newLFUNode(struct{}{}, struct{}{})
	head.next = tail
	tail.prev = head
	return &lfuNodeList{
		head: head,
		tail: tail,
	}
}

func (n *lfuNodeList) removeNode(node *lfuNode) {
	if node == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	if node.prev == nil {
		node.next = nil
		return
	}

	if node.next == nil {
		node.prev.next = nil
		return
	}
	node.prev.next, node.next.prev = node.next, node.prev
}

func (n *lfuNodeList) addNode(node *lfuNode) {

	if node == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()
	node.next = n.head.next
	if n.head.next == nil {
		return
	}
	n.head.next.prev = node
	n.head.next = node
	node.prev = n.head
}
