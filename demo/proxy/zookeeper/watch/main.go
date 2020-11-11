package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/captainlee1024/gateway-demo/proxy/zookeeper"
)

func main() {
	// 获取 zk 节点列表

	zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	// 连接到注册中心（ｚookeeper）的时候首先获取一遍节点列表
	zlist, err := zkManager.GetServerListByPath("/real_server")
	// 拿到之后就可以把这个 zlist 放到负载均衡器里使用了
	fmt.Println("server node:")
	fmt.Println(zlist)
	if err != nil {
		log.Println(err)
	}

	// 动态监听节点变化
	// chanList 接收节点所有内容 chanErr 接收错误
	chanList, chanErr := zkManager.WatchServerListByPath("/real_server")
	go func() {
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr:", changeErr)
			case changedList := <-chanList:
				// 把变化之后的列表注册到负载均衡器里
				fmt.Println("watch node changed:")
				fmt.Println(changedList)
			}
		}
	}()

	/*
		// 获取节点内容
		// 连接之后先做内容同步处理
		zc, _, err := zkManager.GetPathData("/rs_server_conf")
		if err != nil {
			log.Println(err)
		}
		fmt.Println("get node data:")
		fmt.Println(string(zc))

		// 动态监听节点内容
		// 然后异步监听去同步变动
		dataChan, dataErrChan := zkManager.WatchPathData("/rs_server_conf")
		go func() {
			for {
				select {
				case changeErr := <-dataErrChan:
					fmt.Println("changeErr")
					fmt.Println(changeErr)
				case cahngedData := <-dataChan:
					fmt.Println("WatchGetData changed")
					fmt.Println(string(cahngedData))
				}
			}
		}()
	*/

	// 关闭信号监听
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
