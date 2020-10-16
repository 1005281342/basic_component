package load_balance

type LoadBalance interface {
	// Add 添加节点
	Add(...string) error

	// Get 获取节点
	Get(string) (string, error)

	// Update
	Update()
}
