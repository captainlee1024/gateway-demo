package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	addr = "127.0.0.1:2002"
)

func main() {
	// 127.0.0.1:2002/xxx
	// 127.0.0.1:2003/base/xxx
	// 1. 首先拿到目标服务器的地址
	rs1 := "http://127.0.0.1:2003/base"
	// 2. 解析出要代理的url
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}

	// 3. 使用要代理的url创建处ReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(url1)
	log.Println("starting httpserver at " + addr)
	// 4. 开启代理服务器，并代理目标服务器
	// 把proxy当做路由处理器来使用，因为ReverseProxy实现了ServeHTTP方法
	log.Fatal(http.ListenAndServe(addr, proxy))
}
