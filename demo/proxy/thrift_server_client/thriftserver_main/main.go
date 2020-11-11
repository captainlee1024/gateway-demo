package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/captainlee1024/gateway-demo/demo/proxy/thrift_server_client/gen-go/thrift_gen"
)

// Addr 服务器端口
const Addr = "127.0.0.1:6001"

// FormatDataImpl 实现 thrift 里的结构体方法
type FormatDataImpl struct{}

// DoFormat 实现
func (fdi *FormatDataImpl) DoFormat(ctx context.Context, data *thrift_gen.Data) (r *thrift_gen.Data, err error) {
	var rData thrift_gen.Data
	rData.Text = Addr + " DoFormat from server"
	return &rData, nil
}

func main() {
	addr := flag.String("addr", Addr, "input addr")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	handler := &FormatDataImpl{}
	processor := thrift_gen.NewFormatDataProcessor(handler)
	serverSocket, err := thrift.NewTServerSocket(*addr)
	if err != nil {
		log.Fatal("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverSocket, transportFactory, protocolFactory)
	fmt.Println("Running at:", *addr)
	if err := server.Serve(); err != nil {
		log.Fatal(err.Error())
	}
}
