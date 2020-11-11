package loadbalance

import (
	"fmt"

	"github.com/captainlee1024/gateway-demo/proxy/zookeeper"
)

// Conf 服务发现配置主题
type Conf interface {
	Attach(o Observer)        // 绑定
	GetConf() []string        // 获取配置
	WatchConf()               // 监听
	UpdateConf(conf []string) // 更新
}

// ZKConf zookeeper 客户端服务发现配置
type ZKConf struct {
	Observers    []Observer        // 观察者列表
	path         string            // 当前服务器的 zk 的地址
	zkHosts      []string          // zk Hosts 集群的 IP 列表
	confIPWeight map[string]string // 服务器配置的权重和 IP 的 map 结构
	activeList   []string          // 当前可用的 IP 列表
	format       string            // format 格式
}

// NotifyAllObservers 通知所有的 Observer
func (s *ZKConf) NotifyAllObservers() {
	// 被通知的每个观察者会调用里面的 Update() 方法
	for _, obs := range s.Observers {
		obs.Update()
	}
}

// Attach 把观察者绑定到 Observer 上面
func (s *ZKConf) Attach(o Observer) {
	s.Observers = append(s.Observers, o)
}

// GetConf 获取当前 ZK 注册中心的服务列表
func (s *ZKConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIPWeight[ip]
		if !ok {
			weight = "50" // 默认权重
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// WatchConf 更新配置，通知监听者也更新
func (s *ZKConf) WatchConf() {
	zkManager := zookeeper.NewZkManager(s.zkHosts)
	zkManager.GetConnect()
	fmt.Print("watchConf")
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)
	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr", changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed")
				s.UpdateConf(changedList)
			}
		}
	}()
}

// UpdateConf 更新配置时，通知监听者也更新
func (s *ZKConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.Observers {
		obs.Update()
	}
}

// NewLoadBalanceZkConf 创建 ZKConf
func NewLoadBalanceZkConf(
	format, path string, zkHosts []string, conf map[string]string) (*ZKConf, error) {
	zkManager := zookeeper.NewZkManager(zkHosts)
	zkManager.GetConnect()
	defer zkManager.Close()
	zlsit, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}
	mConf := &ZKConf{
		format:       format,
		activeList:   zlsit,
		confIPWeight: conf,
		zkHosts:      zkHosts,
		path:         path,
	}
	mConf.WatchConf()
	return mConf, nil
}

// Observer 观察者接口
type Observer interface {
	Update()
}

// LoadBalanceObserver 观察者
type LoadBalanceObserver struct {
	ModuleConf *ZKConf
}

// Update 更新配置
func (l *LoadBalanceObserver) Update() {
	fmt.Println("Update get conf:", l.ModuleConf.GetConf())
}

// NewLoadBalanceObserver 创建 NewLoadBalanceObserver
func NewLoadBalanceObserver(conf *ZKConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}
