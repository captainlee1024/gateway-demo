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
	rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
	rb.Add("http://127.0.0.1:2003", "50")
	rb.Add("http://127.0.0.1:2004", "100")
	proxy := proxy2.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
