package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/captainlee1024/gateway-demo/proxy/public"

	"github.com/captainlee1024/gateway-demo/proxy/proxy"

	"github.com/captainlee1024/gateway-demo/proxy/middleware"
)

var addr = "127.0.0.1:2002"

// 熔断方案
func main() {
	// 核心方法
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"
		rs2 := "http://127.0.0.1:2004/base"
		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			log.Println(err1)
		}
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			log.Println(err2)
		}

		urls := []*url.URL{url1, url2}
		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}
	log.Println("Starting httpserver at " + addr)

	public.ConfCircuitBreaker(true)

	// 初始化方法数组路由
	sliceRouter := middleware.NewSliceRouter()
	// 路由组的方法数组中添加中间件方法
	sliceRouter.Group("/").Use(middleware.CircuitMW())
	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))

}
