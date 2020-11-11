package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. 建立连接，拿到一个socket连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})
	if err != nil {
		fmt.Printf("conn failed, err: %v\n", err)
		return
	}

	for i := 0; i < 100; i++ {
		// 2. 发送数据
		// 直接在这个套接字上write就行
		_, err := conn.Write([]byte("hello server!"))
		if err != nil {
			fmt.Printf("send data failed, err: %v\n", err)
			return
		}

		// 3. 接收数据
		result := make([]byte, 1024)
		// 接收数据也是直接在这个套接字上进行read就行
		n, remoteAddr, err := conn.ReadFromUDP(result)
		if err != nil {
			fmt.Printf("receive failed, err: %v\n", err)
			return
		}
		fmt.Printf("==> receive from affr: %v data: %v\n", remoteAddr, string(result[:n]))
	}
}
