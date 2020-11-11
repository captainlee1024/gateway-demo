package loadbalance

// LoadBalance 负载均衡策略接口
type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	// 后期服务发现
	Update()
}
