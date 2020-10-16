package hashring

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sync"
)

// 哈希环
type HashRing struct {
	nodes        sync.Map  // key: index, value: node
	idxList      *SkipList // 基于跳表存储
	replicaCount int       // 虚拟节点数
	hash         hash.Hash32
}

// 新建一个HashRing，设置虚拟节点、Hash函数，如果hash函数为nil，则使用fnv32a
func New(replicaCount int, hash hash.Hash32) *HashRing {
	if hash == nil {
		hash = fnv.New32a()
	}

	return &HashRing{
		nodes:        sync.Map{},
		replicaCount: replicaCount,
		hash:         hash,
		idxList:      newSkipList(),
	}
}

// 计算HashCode
func getHashCode(hash hash.Hash32, key []byte) (uint32, error) {
	hash.Reset()
	_, err := hash.Write(key)
	if err != nil {
		return 0, err
	}

	return hash.Sum32(), nil
}

// 添加一个节点到HashRing中
func (hr *HashRing) Add(node string) error {

	for i := 0; i < hr.replicaCount; i++ {
		key := fmt.Sprintf("%s:%d", node, i)
		hKey, err := getHashCode(hr.hash, []byte(key))
		if err != nil {
			return fmt.Errorf("failed to add node: %v", err)
		}

		// 选择的跳表、Map都是线程安全的实现，不需要加锁
		hr.idxList.Set(float64(hKey), node)
		hr.nodes.Store(hKey, node)
	}

	return nil
}

// 从环中移除节点
func (hr *HashRing) Delete(node string) error {

	for i := 0; i < hr.replicaCount; i++ {
		key := fmt.Sprintf("%s:%d", node, i)
		hKey, err := getHashCode(hr.hash, []byte(key))
		if err != nil {
			return fmt.Errorf("failed to delete node: %v", err)
		}

		// 移除节点
		hr.nodes.Delete(hKey)
		hr.idxList.Remove(float64(hKey))
	}
	return nil
}

// Locate returns the node for a given key
func (hr *HashRing) Locate(key string) (node string, err error) {

	if hr.idxList.Length < 1 {
		return node, fmt.Errorf("no available nodes")
	}

	var hKey uint32
	if hKey, err = getHashCode(hr.hash, []byte(key)); err != nil {
		return node, fmt.Errorf("failed to fetch node: %v", err)
	}

	return hr.idxList.Get(float64(hKey)).Value().(string), nil
}
