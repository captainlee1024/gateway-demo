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
	mConf, err := loadbalance.NewCheckConf(
		"http://%s/base",
		map[string]string{"127.0.0.1:2003": "20", "127.0.0.1:2004": "20"})
	// nil)
	if err != nil {
		panic(err)
	}

	rb := loadbalance.FactorWithConf(loadbalance.LbWeightRoundRobin, mConf)
	proxy := proxy2.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
