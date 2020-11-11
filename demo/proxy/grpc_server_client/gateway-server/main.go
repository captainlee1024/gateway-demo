package main

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"flag"
	"fmt"
	"net/http"

	gw "github.com/captainlee1024/gateway-demo/demo/proxy/grpc_server_client/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/golang/glog"
)

var (
	serverAddr         = ":8081"
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50055", "gRPC server endpoint")
)

func run() error {
	// 创建上下文
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建一个类似http路由器的mux，这里用来跟下游建立连接
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterEchoHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(serverAddr, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	fmt.Println("server listening at", serverAddr)
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
