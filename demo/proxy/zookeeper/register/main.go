package main

import (
	"fmt"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/zookeeper"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	i := 0
	for {
		// 每 5 秒注册一个
		zkManager.RegistServerPath("/real_server", fmt.Sprint(i))
		fmt.Println("Register", i)
		time.Sleep(time.Second * 5)
		i++
	}
}
