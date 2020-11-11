package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var (
	proxyAddr1 = "http://127.0.0.1:2003"
	proxyAddr2 = "http://127.0.0.1:2004"
	port       = "2002"
	count      = 0
)

func reversePxy(w http.ResponseWriter, r *http.Request) {
	// 1. 解析代理地址，更新请求体的协议和主机（更改请求结构体信息）
	count++
	if (count % 2) == 0 {
		proxy, err := url.Parse(proxyAddr1)
		if err != nil {
			fmt.Printf("Parse prox_addr failed, err: %v\n", err)
			return
		}
		r.URL.Scheme = proxy.Scheme
		r.URL.Host = proxy.Host

	} else {
		proxy, err := url.Parse(proxyAddr2)
		if err != nil {
			fmt.Printf("Parse prox_addr failed, err: %v\n", err)
			return
		}
		r.URL.Scheme = proxy.Scheme
		r.URL.Host = proxy.Host

	}

	// 2. 请求下游（通过负载均衡获取下有服务器地址，并把更改过的请求发送到服务器）
	// 创建默认transport
	transport := http.DefaultTransport
	// 请求下游处理方法
	resp, err := transport.RoundTrip(r)
	if err != nil {
		log.Print(err)
		return
	}

	// 3. 对返回的信息做一些处理（拷贝header和响应内容）然后返回给客户端。
	// 拷贝请求头Header信息
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()
	// 拷贝请求内容
	bufio.NewReader(resp.Body).WriteTo(w)
}

func main() {
	http.HandleFunc("/", reversePxy)
	log.Println("Start serving on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
