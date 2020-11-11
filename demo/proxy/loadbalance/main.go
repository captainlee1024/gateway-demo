package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
)

var (
	addr      = "127.0.0.1:2002"
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 30, // 连接超时时间
			KeepAlive: time.Second * 30, // 长连接超时时间
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接
		IdleConnTimeout:       time.Second * 90, // 空闲超时时间
		TLSHandshakeTimeout:   time.Second * 10, // tls 握手超时时间
		ExpectContinueTimeout: time.Second * 1,  // 100-continue 状态码超时时间
	}
)

// NewMultipleHostsReverseProxy 反向代理负载均衡
func NewMultipleHostsReverseProxy(lb loadbalance.LoadBalance) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.RemoteAddr)
		// nextAddr, err := lb.Get(req.URL.String())
		if err != nil {
			log.Fatal("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoinSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		// 请求以下命令：curl 'http://127.0.0.1:2002/error'
		if resp.StatusCode != 200 {
			// 获取内容
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			// 追加内容
			newPayload := []byte("StatusCode error:" + string(oldPayload))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			resp.ContentLength = int64(len(newPayload))
			resp.Header.Set("Content", strconv.FormatInt(int64(len(newPayload)), 10))
		}
		return nil
	}

	// 错误回调：关闭real_server时测试，错误回调
	// 范围：transport.RoundTrip 发生错误，以及 ModifyResponse 发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		// todo 如果是权重的负载则调整临时权重
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}
	return &httputil.ReverseProxy{
		Director:       director,   // 请求协调者
		Transport:      transport,  //　连接池
		ModifyResponse: modifyFunc, // 修改返回内容
		ErrorHandler:   errFunc,    // 错误回调函数
	}

}

func singleJoinSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func main() {
	rb := loadbalance.Factory(loadbalance.LbConsistentHash)
	if err := rb.Add("http://127.0.0.1:2003/bse", "10"); err != nil {
		log.Println(err)
	}
	if err := rb.Add("http://127.0.0.1:2004/base", "20"); err != nil {
		log.Println(err)
	}
	proxy := NewMultipleHostsReverseProxy(rb)
	log.Println("Starting httpserver at ", addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
