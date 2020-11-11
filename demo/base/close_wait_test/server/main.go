package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. 监听端口
	// 这里监听本地的 9090 端口，拿到监听的 listener
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}

	// 2. 建立套接字连接
	for {
		// 用循环持续的去获取监听的端口的 connection 连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			return
		}
		// 3. 创建处理协程
		// 拿到 connection 连接之后，创建协程，使用 connection 持续的去读取信息
		// 这里一定要 defer 调用关闭
		// 如果这里不关闭服务端会有CLOSE-WAIT状态，客户端会有FIN-WAIT状态
		// 每一次读取128字节长度的信息，当出现error的时候跳出循环，协程结束
		// 当客户端关闭的时候会出现error:EOF，这个时候退出
		go func(conn net.Conn) {
			defer conn.Close() // 如果这里不写会有什么问题
			for {
				var buf [128]byte
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Printf("read from connect failed, err: %v\n", err)
					break
				}
				str := string(buf[:n])
				fmt.Printf("receive from client, data: %v\n", str)
			}
		}(conn)
	}
}
