package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}
	// 2. 建立套接字连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed, err: %v\n", err)
			continue
		}

		// 3. 创建协程
		go process(conn)
	}
}

func process(conn net.Conn) {

	/*
		如果不关闭，在四次挥手的时候，被动关闭方不会进行 CLOSE-WAIT，一直处于 CLOSE-WAIT 状态
		而且也不会发送 FIN 包到主动关闭方，所以会一直处于FIN-WAIT-2 状态
	*/
	defer conn.Close() // 如果不关闭会出现什么问题
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
}
