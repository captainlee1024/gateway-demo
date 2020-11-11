package loadbalance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// WeightRoundRobinBalance 加权轮询
type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	// 观察主题
	conf Conf
}

// WeightNode 权重节点
type WeightNode struct {
	addr            string
	Weight          int // 权重
	currentWeight   int // 临时权重
	effectiveWeight int // 有效权重
}

// Add 添加下游服务
// params 接收两个string参数 第一个为服务器Ip 第二个是该服务节点权重值
func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{
		addr:   params[0],
		Weight: int(parInt),
	}
	node.effectiveWeight = node.Weight
	r.rss = append(r.rss, node)
	return nil
}

// Next 加权获取下游服务
func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]

		// 1. 统计所有有效权重之和 sum(EffectiveWeight)
		total += w.effectiveWeight

		// 2. 变更节点临时权重＝节点临时权重＋节点有效权重
		w.currentWeight += w.effectiveWeight

		// 3. 有效权重默认与权重相同，通讯异常时 -1，通讯成功 +1，直到恢复到 weight 大小
		if w.effectiveWeight < w.Weight {
			w.effectiveWeight++
		}

		// 4. 选择最大临时权重节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}

	}
	if best == nil {
		return ""
	}
	// 5. 变更临时权重为临时权重－有效权重之和
	best.currentWeight -= total
	return best.addr
}

// Get Next方法的进一步包装
func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

// SetConf 设置服务发现
func (r *WeightRoundRobinBalance) SetConf(conf Conf) {
	r.conf = conf
}

// Update 更新，通知其他...
func (r *WeightRoundRobinBalance) Update() {
	if conf, ok := r.conf.(*ZKConf); ok {
		fmt.Println("WeightRoundRobinBalance get conf:", conf.GetConf())
		r.rss = nil
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*CheckConf); ok {
		fmt.Println("WeightRoundRobinBalance get conf:", conf.GetConf())
		r.rss = nil
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
