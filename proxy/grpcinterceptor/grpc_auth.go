package grpcinterceptor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// valid validates the authorization.
func valid(authorization []string) bool {
	fmt.Println("len(authrization)", len(authorization))
	// 长度小于1，直接返回false
	if len(authorization) < 1 {
		return false
	}
	// 长度不小于１，验证前缀是不是 "Bearer " 内容是不是 some-select-token
	// 如果是就放行，返回true
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}

// wrappedStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RcvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// GrpcAuthStreamInterceptor 流式 RPC 拦截器
func GrpcAuthStreamInterceptor(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 通过 ServerStream 取出上下文，通过上下文拿到服务的 metadata
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errMissingMetadata
	}

	// 获取 md 中的 authorization Header 头，并校验该 Header 头的长度
	if !valid(md["authorization"]) {
		return errInvalidToken
	}

	// 请求下游服务
	err := handler(srv, ss)
	if err != nil {
		log.Printf("RPC failed with error %v\n", err)
	}
	return err
}

// GrpcAuthUnaryInterCeptor 普通 RPC 拦截器
func GrpcAuthUnaryInterCeptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("GrpcAuthUnaryInterceptor")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error%v\n", err)
	}
	return m, err
}
