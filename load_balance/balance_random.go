package load_balance

import "math/rand"

type BalanceWithRandom struct {
	// 均衡器公共功能
	*balance

	// 节点列表
	rss []string
}

func NewBalanceWithRandom() *BalanceWithRandom {
	return &BalanceWithRandom{rss: make([]string, 0)}
}

func (b *BalanceWithRandom) Add(params ...string) error {
	if len(params) <= 0 {
		return ErrParamLeastOne
	}
	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	// 将节点添加到节点列表
	b.rss = append(b.rss, params[0])
	return nil
}

func (b *BalanceWithRandom) Get(key string) (string, error) {
	return b.getNode()
}

func (b *BalanceWithRandom) getNode() (string, error) {
	if len(b.rss) == 0 {
		return "", ErrNodeListIsEmpty
	}

	b.rwLock.RLock()
	defer b.rwLock.RUnlock()

	// 随机获取节点
	return b.rss[rand.Intn(len(b.rss))], nil
}

func (b *BalanceWithRandom) Update() {}
