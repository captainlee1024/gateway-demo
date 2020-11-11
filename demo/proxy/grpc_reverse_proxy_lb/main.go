package main

import (
	"fmt"
	"log"
	"net"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
	proxy2 "github.com/captainlee1024/gateway-demo/proxy/proxy"
	"github.com/captainlee1024/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	rb := loadbalance.Factory(loadbalance.LbWeightRoundRobin)
	rb.Add("127.0.0.1:50055", "40")

	grpcHandler := proxy2.NewGrpcLoadBalanceHandler(rb)
	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(grpcHandler)) // 自定义全局回调

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v\n", err)
	}
}
