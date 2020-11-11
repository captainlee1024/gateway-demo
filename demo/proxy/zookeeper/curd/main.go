package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	host = []string{"127.0.0.1:2181"}
)

func main() {
	// 连接服务器
	// 设置 ip 、超时时间
	conn, _, err := zk.Connect(host, time.Second*5)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 增
	// 直接使用拿到的连接句柄创建
	// path "/test_tree2"，在 path 中设置 byte 数组
	// 0 是节点类型
	// 最后一个是权限，这里设置成大家都可以进行增删改查 zk.WorldACL(zk.PermAll)
	if _, err := conn.Create("/test_tree2", []byte("tree_connect"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	// 查
	// 直接使用句柄调用 Get 方法，拿到的是 []byte 二进制数组
	nodeValue, dStat, err := conn.Get("/test_tree2")
	if err != nil {
		fmt.Println("get err", err)
		return
	}
	fmt.Println("nodeValue", string(nodeValue))

	// 改
	// 修改要基于版本号进行修改，所以修改要先进行查询获取版本号
	if _, err := conn.Set("/test_tree2", []byte("new_content"), dStat.Version); err != nil {
		fmt.Println("udpate err", err)
	}

	// // 查
	// // 直接使用句柄调用 Get 方法，拿到的是 []byte 二进制数组
	// nodeValue, dStat, err = conn.Get("/test_tree2")
	// if err != nil {
	// 	fmt.Println("get err", err)
	// 	return
	// }
	// fmt.Println("nodeValue", string(nodeValue))

	// 删除
	// 删除也需要版本号，所以也需要先查询一下
	_, dStat, _ = conn.Get("/test_tree2")
	if err := conn.Delete("/test_tree2", dStat.Version); err != nil {
		fmt.Println("Delete err", err)
		// return
	}

	// 验证存在
	hasNode, _, err := conn.Exists("/test_tree2")
	if err != nil {
		fmt.Println("Exists err", err)
		// return
	}
	fmt.Println("node Exist", hasNode)

	// 增加
	if _, err := conn.Create("/test_tree2", []byte("tree_content"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	// 设置子节点
	// 设置子节点时，上游节点不存在会报错
	if _, err := conn.Create("/test_tree2/subnode", []byte("node_content"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	// 获取子节点列表
	childNodes, _, err := conn.Children("/test_tree2")
	if err != nil {
		fmt.Println("children err", err)
	}
	fmt.Println("children", childNodes)

	// 删除
	// 删除也需要版本号，所以也需要先查询一下
	// _, dStat, _ = conn.Get("/test_tree2")
	// if err := conn.Delete("/test_tree2", dStat.Version); err != nil {
	// 	fmt.Println("Delete err", err)
	// 	// return
	// }
}
