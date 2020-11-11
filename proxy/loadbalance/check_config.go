package loadbalance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

//default check setting
const (
	// DefaultCheckMethod xx
	DefaultCheckMethod = 0
	// DefaultCheckTimeout xx
	DefaultCheckTimeout = 2
	// DefaultCheckMaxErrNum xx
	DefaultCheckMaxErrNum = 2
	// DefaultCheckInterval xx
	DefaultCheckInterval = 5
)

// CheckConf 客户端服务发现负载均衡配置
type CheckConf struct {
	observers    []Observer
	confIPWeight map[string]string
	activeList   []string
	format       string
}

// Attach 绑定
func (s *CheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers 通知所有的
func (s *CheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf 获取配置
func (s *CheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIPWeight[ip]
		if !ok {
			weight = "50" //默认weight
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// WatchConf 更新配置时，通知监听者也更新
func (s *CheckConf) WatchConf() {
	fmt.Println("watchConf")
	go func() {
		// 记录每个 IP 节点的错误数的 map
		confIPErrNum := map[string]int{}
		for {
			// 存储发现的所有的可用服务节点
			changedList := []string{}
			// for item, _ := range s.confIPWeight {
			// 遍历所有服务节点进行心跳探测（Dial 拨号探测）
			for item := range s.confIPWeight {
				// 每拿到一个下游节点就进行一个 Dial(尝试三次握手)
				// 如果没能在 DefaultCheckTimeout 拨号成功就记录一个错误信息
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				//TODO http statuscode的检测，还可以做 rpc 的服务接口的检测
				// 没有出错就设置该 Ip 在 confIPErrNum map 中对应的值设为０
				if err == nil {
					conn.Close()
					if _, ok := confIPErrNum[item]; ok {
						confIPErrNum[item] = 0
					}
				}
				// 出现错误设置该 Ip 在 confIPErrNum map 中对应的值累计加 1
				if err != nil {
					if _, ok := confIPErrNum[item]; ok {
						// confIPErrNum[item] += 1
						confIPErrNum[item]++

					} else {
						confIPErrNum[item] = 1
					}
				}
				// 验证错误数
				// 如果错误信息数没有超过最大的错误数，就放到可用的服务节点列表中
				if confIPErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			// 把 changedList 排序
			sort.Strings(changedList)
			// 当前可用节点排序
			sort.Strings(s.activeList)
			// 比较之前的可用节点列表 s.activeList 和本次探测处的可用节点 changedList 是否一致
			// 不一致就进行更新
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			// 休眠固定的时间进行下一轮探测
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

// UpdateConf 更新配置时，通知监听者也更新
func (s *CheckConf) UpdateConf(conf []string) {
	fmt.Println("UpdateConf", conf)
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

// NewCheckConf 创建主动探测的负载均衡配置 CheckConf
func NewCheckConf(format string, conf map[string]string) (*CheckConf, error) {
	aList := []string{}
	//默认初始化
	for item := range conf {
		aList = append(aList, item)
	}
	mConf := &CheckConf{format: format, activeList: aList, confIPWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}
