package main

import (
	"fmt"
	"net"

	"github.com/captainlee1024/gateway-demo/demo/base/unpack/unpack"
)

func main() {
	// 首先连接服务器拿到一个套接字
	conn, err := net.Dial("tcp", "localhost:9090")
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err.Error())
		return
	}

	// 通过套接字执行Encode方法，写入想要发送的数据到缓冲区，TCP发送到接收缓冲区，另一方再读取
	unpack.Encode(conn, "hello world 呀0!!!")
}
