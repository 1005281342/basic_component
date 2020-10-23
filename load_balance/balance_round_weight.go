package load_balance

import (
	"strconv"
)

type BalanceWithRoundWeight struct {
	*balance
	// 带权重节点列表
	rss []*weightNode
}

type weightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func NewBalanceWithRoundWeight() *BalanceWithRoundWeight {
	return &BalanceWithRoundWeight{rss: make([]*weightNode, 0)}
}

func (b *BalanceWithRoundWeight) Add(params ...string) error {
	if len(params) < 2 {
		return ErrParamLeastTwo
	}

	var (
		parInt int64
		err    error
	)

	if parInt, err = strconv.ParseInt(params[1], 10, 64); err != nil {
		return err
	}

	// 初始化权重节点
	var node = &weightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight

	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	// 将权重节点添加到权重节点列表
	b.rss = append(b.rss, node)
	return nil
}

func (b *BalanceWithRoundWeight) Get(key string) (string, error) {
	return b.getNode()
}

func (b *BalanceWithRoundWeight) getNode() (string, error) {
	var (
		total int
		best  *weightNode
		w     *weightNode
	)

	b.rwLock.RLock()
	defer b.rwLock.RUnlock()
	for i := 0; i < len(b.rss); i++ {
		w = b.rss[i]
		//step 1 统计所有有效权重之和
		total += w.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}

	// 没有可用节点
	if best == nil {
		return "", ErrNodeNotAvailable
	}

	//step 5 变更临时权重为 临时权重-有效权重之和
	best.currentWeight -= total
	return best.addr, nil
}

func (b *BalanceWithRoundWeight) Update() {}
