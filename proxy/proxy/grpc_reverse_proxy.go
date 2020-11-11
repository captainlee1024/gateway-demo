package proxy

import (
	"context"
	"log"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
	"github.com/captainlee1024/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

// NewGrpcLoadBalanceHandler description
/*
func NewGrpcLoadBalanceHandler(lb loadbalance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			// 拒绝某些特殊请求
			if strings.HasPrefix(fullMethodName, "/come.example.internal.") {
				return ctx, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
			}
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				// 基于 header 头决定下游请求
				if val, exists := md[":authority"]; exists && val[0] == "staging.api.example.com" {
					// return ctx, nil, nil
					// TODO 权限认证，待完善，不清处此时应该怎么处理
					return ctx, nil, grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
				}
				c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
				return ctx, c, err
			}
			// c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			// return ctx, c, err
			return ctx, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
		}
		return proxy.TransparentHandler(director)
	}()
}
*/

// NewGrpcLoadBalanceHandler 创建 Grpc 负载均衡代理 Handler
func NewGrpcLoadBalanceHandler(lb loadbalance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			return ctx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
