package main

import (
	"log"
	"net/http"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
	"github.com/captainlee1024/gateway-demo/proxy/middleware"
	proxy2 "github.com/captainlee1024/gateway-demo/proxy/proxy"
)

var (
	addr = "127.0.0.1:2002"
)

func main() {

	// 创建 loadbalanceZKConf
	mConf, err := loadbalance.NewLoadBalanceZkConf(
		"http://%s/base", "/real_server", []string{"127.0.0.1:2181"},
		map[string]string{"127.0.0.1:2003": "20"})
	if err != nil {
		panic(err)
	}

	// 利用 loadbalanceZKConf 构建负载均衡器
	// rb := loadbalance.FactorWithConf(loadbalance.LbWeightRoundRobin, mConf)
	rb := loadbalance.FactorWithConf(loadbalance.LbConsistentHash, mConf)

	// 利用负载均衡器构建方向代理服务器
	proxy := proxy2.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
	log.Println("Starting httpserver at" + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
