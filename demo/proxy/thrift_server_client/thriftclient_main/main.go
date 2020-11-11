package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/captainlee1024/gateway-demo/demo/proxy/thrift_server_client/gen-go/thrift_gen"

	"github.com/apache/thrift/lib/go/thrift"
)

func main() {
	// 直接访问thrift服务
	// addr := flag.String("addr", "127.0.0.1:6001", "input addr")

	// 访问tcp代理服务器（用来代理thrift服务）
	addr := flag.String("addr", "127.0.0.1:2002", "input addr")

	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	for {
		tSocket, err := thrift.NewTSocket(*addr)
		if err != nil {
			log.Fatal("tSocket error:", err)
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, _ := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thrift_gen.NewFormatDataClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Fatal("Error opening:", *addr)
		}
		defer transport.Close()
		data := thrift_gen.Data{Text: "ping"}
		d, err := client.DoFormat(context.Background(), &data)
		if err != nil {
			fmt.Println("err:", err.Error())
		} else {
			fmt.Println("Text:", d.Text)
		}
		time.Sleep(time.Millisecond * 40)
	}
}
