package main

import (
	"fmt"
	"net"
)

func main() {
	// UDP服务器读取消息不需要监听后再建立socket连接使用套接字去读取
	// 而是直接使用监听的listen句柄去读取信息
	// 1. 监听服务器
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})
	if err != nil {
		fmt.Printf("listen failed, err: %v\n", err)
		return
	}

	// 2. 循环读取消息内容
	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			fmt.Printf("read failed from addr: %v, err: %v\n", addr, err)
			break
		}

		go func() {
			// todo sth
			// 3. 回复数据
			// 先打印一遍获取数据的信息
			fmt.Printf("addr: %v data: %v count: %v\n", addr, string(data[:n]), n)
			// 写入恢复的数据
			_, err := listen.WriteToUDP([]byte("received success!"), addr)
			if err != nil {
				fmt.Printf("write failed, err: %v\n", err)
			}
		}()
	}
}
