package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	// 1. 创建连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 30, // 连接超时
			KeepAlive: time.Second * 30, // 探活时间
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接
		IdleConnTimeout:       time.Second * 90, // 空闲超时时间
		TLSHandshakeTimeout:   time.Second * 1,  // tls 握手超时时间
		ExpectContinueTimeout: time.Second * 1,  // 100-continue状态码超时时间
	}
	// 2. 创建客户端
	client := &http.Client{
		Timeout:   time.Second * 30, // 请求超时时间
		Transport: transport,
	}

	// 3. 请求数据
	resp, err := client.Get("http://127.0.0.1:1210/bye")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 4. 读取内容
	bds, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bds))
}
