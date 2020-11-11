package main

import (
	"fmt"

	"github.com/captainlee1024/gateway-demo/proxy/grpcinterceptor"
	proxy2 "github.com/captainlee1024/gateway-demo/proxy/proxy"
	"github.com/captainlee1024/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"log"
	"net"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
	"github.com/captainlee1024/gateway-demo/proxy/public"
)

const port = ":50051"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
	rb.Add("127.0.0.1:50055", "40")

	counter, _ := public.NewFlowCountService("local_app", time.Second)
	grpcHandler := proxy2.NewGrpcLoadBalanceHandler(rb)
	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpcinterceptor.GrpcAuthStreamInterceptor,
			grpcinterceptor.GrpcFlowCountStreamInterceptor(counter)), // 流式方法拦截
		grpc.CustomCodec(proxy.Codec()),         // 自定义 codec
		grpc.UnknownServiceHandler(grpcHandler)) // 自定义全局回调

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v\n", err)
	}
}
