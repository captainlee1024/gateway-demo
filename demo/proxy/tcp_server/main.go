package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/captainlee1024/gateway-demo/proxy/tcpproxy"
)

var (
	addr = ":7002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func main() {

	// tcp 服务器测试
	log.Println("Starting tcpserver at " + addr)
	tcpServer := tcpproxy.TCPServer{
		Addr:    addr,
		Handler: &tcpHandler{},
	}
	fmt.Println("Starting tcpserver at " + addr)
	tcpServer.ListenAndServe()

	/*
		// tcp 代理测试
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:6001", "40")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: proxy,
		}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServer.ListenAndServe()
	*/
}
