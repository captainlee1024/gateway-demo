package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/captainlee1024/gateway-demo/proxy/proxy"

	"github.com/captainlee1024/gateway-demo/proxy/middleware"
)

var addr = "127.0.0.1:2002"

// 限流方案
func main() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		// 两个下游服务
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

	// 构造方法数组路由器
	sliceRouter := middleware.NewSliceRouter()
	// 创建路由组，并添加中间件到路由组的方法切片中
	sliceRouter.Group("/").Use(middleware.RateLimiter())
	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}
