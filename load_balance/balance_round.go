package load_balance

type BalanceWithRound struct {
	*BalanceWithRandom
	// 当前节点下标
	curIndex int
}

func NewBalanceWithRound() *BalanceWithRound {
	return &BalanceWithRound{}
}

func (b *BalanceWithRound) Get(key string) (string, error) {
	return b.getNode()
}

func (b *BalanceWithRound) getNode() (string, error) {
	if len(b.rss) == 0 {
		return "", ErrNodeListIsEmpty
	}

	b.rwLock.RLock()
	defer b.rwLock.RUnlock()

	// 循环获取节点
	var lens = len(b.rss)
	if b.curIndex >= lens {
		b.curIndex = 0
	}
	var curAddr = b.rss[b.curIndex]
	b.curIndex = (b.curIndex + 1) % lens
	return curAddr, nil
}
