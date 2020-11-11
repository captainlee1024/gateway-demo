package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/proxy"
	"github.com/captainlee1024/gateway-demo/proxy/tcpproxy"

	"github.com/captainlee1024/gateway-demo/proxy/tcpmiddleware"

	"github.com/captainlee1024/gateway-demo/proxy/public"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
)

var (
	addr = ":2002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler"))
}

func main() {

	// 基于 thrift 代理测试 tcp 中间件
	rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
	rb.Add("127.0.0.1:6001", "40")

	// 构建路由并设置中间件
	counter, _ := public.NewFlowCountService("local_app", time.Second)
	router := tcpmiddleware.NewTCPSliceRouter()
	router.Group("/").Use(tcpmiddleware.IPWhiteListMiddleWare(),
		tcpmiddleware.FlowCountMiddleWare(counter))

	// 构建回调 handler
	routerHandler := tcpmiddleware.NewTCPSliceRouterHandler(
		func(c *tcpmiddleware.TCPSliceRouterContext) tcpproxy.TCPHandler {
			return proxy.NewTCPLoadBalanceReverseProxy(c, rb)
		}, router)

	// 启动服务
	tcpServ := tcpproxy.TCPServer{
		Addr:    addr,
		Handler: routerHandler,
	}
	fmt.Println("Starting tcpproxy at " + addr)
	tcpServ.ListenAndServe()

}
