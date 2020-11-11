package main

import (
	"context"
	"net"
)

var (
	addr = ":2002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func main() {
	/*
		// tcp 服务器测试
		log.Println("Starting tcpserver at " + addr)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: &tcpHandler{},
		}
		fmt.Println("Starting tcpserver at " + addr)
		tcpServer.ListenAndServe()
	*/

	/*
		// tcp 代理测试
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:7002", "100")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: proxy,
		}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServer.ListenAndServe()
	*/

	/*
		// thrift 代理测试
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:6001", "100")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: proxy,
		}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServer.ListenAndServe()
	*/

	/*
		// redis 服务测试
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:6379", "40")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: proxy,
		}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServer.ListenAndServe()
	*/

	// http 代理测试
	// 缺点对请求的管控不足，比如我们用来做baidu 代理，因为无法更改host，所以很容易把我们拒绝
	/*
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		// rb.Add("127.0.0.1:2003", "40")
		// rb.Add("127.0.0.1:2004", "60")
		rb.Add("www.baidu.com:80", "60")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServer := tcpproxy.TCPServer{
			Addr:    addr,
			Handler: proxy,
		}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServer.ListenAndServe()
	*/

	//websocket服务器测试:缺点对请求的管控不足
	/*
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:2003", "40")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServ := tcpproxy.TCPServer{Addr: addr, Handler: proxy}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServ.ListenAndServe()
	*/

	//http2服务器测试:缺点对请求的管控不足
	/*
		rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
		rb.Add("127.0.0.1:3003", "40")
		proxy := proxy.NewTCPLoadBalanceReverseProxy(&tcpmiddleware.TCPSliceRouterContext{}, rb)
		tcpServ := tcpproxy.TCPServer{Addr: addr, Handler: proxy}
		fmt.Println("Starting tcpproxy at " + addr)
		tcpServ.ListenAndServe()
	*/
}
