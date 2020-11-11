package loadbalance

// LbType 加权策略
type LbType int

const (
	// LbRandom 随机负载均衡
	LbRandom LbType = iota
	// LbRoundRobin 轮询负载
	LbRoundRobin
	// LbWeightRoundRobin 加权轮询
	LbWeightRoundRobin
	// LbConsistentHash 一致性hash负载均衡
	LbConsistentHash
)

// Factory 负载策略工厂方法
// 默认使用随机负载均衡策略
func Factory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LbConsistentHash:
		return NewConsistentHashBanlance(10, nil)
	default:
		return &RandomBalance{}
	}
}

// FactorWithConf 观察者模式结合负载均衡器工厂方法
func FactorWithConf(lbType LbType, mConf Conf) LoadBalance {
	//观察者模式
	switch lbType {
	case LbRandom:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbConsistentHash:
		lb := NewConsistentHashBanlance(10, nil)
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbWeightRoundRobin:
		lb := &WeightRoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	default:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	}
}
