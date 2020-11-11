package main

import (
	"fmt"
	"net"

	"github.com/captainlee1024/gateway-demo/demo/base/unpack/unpack"
)

func main() {
	// simple tcp server
	// 1. 监听端口，拿到listener
	// 1. listen ip+port
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}

	// 2. 接收请求
	// 2. accept client request
	for {
		// 从listener中accept一个socket连接(套接字)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed, err: %v\n", err)
			return
		}

		// 3. 创建协程
		// 3. create goroutine for each connect套接字(socket连接)
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		// 把拿到的数据进行解码
		bt, err := unpack.Decode(conn)
		if err != nil { // 如果中间没有出现错误也没有结束(EOF)，就会打印读取的结果
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		// 接收的字节转换成字符串
		str := string(bt)
		// 打印从接收缓冲区读取转换之后的数据
		fmt.Printf("receive from client, data: %v\n", str)
	}
}
