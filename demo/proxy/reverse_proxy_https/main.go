package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/captainlee1024/gateway-demo/demo/proxy/reverse_proxy_https/public"
	"github.com/captainlee1024/gateway-demo/demo/proxy/reverse_proxy_https/testdata"
)

var addr = "example1.com:3002"

func main() {
	rs1 := "https://example1.com:3003"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}
	urls := []*url.URL{url1}
	proxy := public.NewMultipleHostsReverseProxy(urls)
	log.Println("Starting httpserver at " + addr)
	// ListenAndServeTLS 会默认添加上 	http2.ConfigureServer(server, &http2.Server{})
	// 所以这里和 https 那里效果是一样的
	log.Fatal(http.ListenAndServeTLS(addr, testdata.Path("server.crt"), testdata.Path("server.key"), proxy))
}
