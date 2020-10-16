package load_balance

import (
	"github.com/1005281342/basic_component/hashring"
	"hash"
)

type BalanceWithConsistentHash struct {
	*balance

	// 一致性Hash
	hashRing *hashring.HashRing
}

func NewBalanceWithConsistentHash(replicas int, hashFunc hash.Hash32) *BalanceWithConsistentHash {

	return &BalanceWithConsistentHash{hashRing: hashring.New(replicas, hashFunc)}
}

func (b *BalanceWithConsistentHash) Add(params ...string) error {
	if len(params) <= 0 {
		return ErrParamLeastOne
	}

	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	return b.hashRing.Add(params[0])
}

func (b *BalanceWithConsistentHash) Get(key string) (string, error) {
	return b.hashRing.Locate(key)
}

func (b *BalanceWithConsistentHash) Update() {

}
