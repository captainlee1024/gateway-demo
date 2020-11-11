package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/captainlee1024/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

func main() {
	// 创建tcp连接，监听 50051 端口
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen:%v\n", err)
	}
	// 下游请求协调者
	// 请求协调者会注册到 UnknownServiceHandler 里
	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		// 拒绝某些特殊的请求
		if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
			return ctx, nil, status.Errorf(codes.Unimplemented, "Unknown method")
		}
		// 如果不是需要拦截的请求方法，就与下游建立连接
		c, err := grpc.DialContext(ctx, "localhost:50055", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		// outCtx, cancel := context.WithCancel(ctx)
		// defer cancel() 这里不能关闭
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		// 返回 clientconn
		return outCtx, c, err
	}

	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v\n", err)
	}
}
