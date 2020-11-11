package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

	// 3. 调用自己的 NewSingleHostReverseProxy
	proxy := NewSingleHostReverseProxy(url1)
	log.Println("starting httpserver at " + addr)
	// 4. 开启代理服务器，并代理目标服务器
	// 把proxy当做路由处理器来使用，因为ReverseProxy实现了ServeHTTP方法
	log.Fatal(http.ListenAndServe(addr, proxy))
}

// NewSingleHostReverseProxy 基于源码实现修改相应内容
// 新建一个 proxy
// 如果目标 rs 路径是 http://127.0.0.1:2003/base ，
// 请求的路径如果是 http://127.0.0.1:2002/dir ，
// 则实际路径为 http://127.0.0.1:2003/base/dir 。
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	// http://127.0.0.1:2002/dir?name=abc
	// 如果url地址里面有参数，协议，主机地址相关内容，会在这里获取到
	// RawQuery: name=abc
	// Scheme: http
	// Host: 127.0.0.1:2002
	targetQuery := target.RawQuery
	// 创建一个Director类型的函数，实现对请求的代理转发
	director := func(req *http.Request) {
		// 协议的重新定义
		req.URL.Scheme = target.Scheme
		// 主机的重新赋值
		req.URL.Host = target.Host
		// Path 的赋值
		// target.Path: /base
		// req.URL.Path: /dir
		// 合并之后的 req.URL.Path : /base/dir
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		// 参数的合并
		if targetQuery == "" || req.URL.RawQuery == "" {
			// 参数都为空就直接合并
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			// 有一个不为空就使用&进行相加合并
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		// 如果Header中"User-Agent"没有值，就设置"User-Agent"，值为""
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	modifyFunc := func(res *http.Response) error {
		// 判断，在服务器返回响应的状态码不是 200 的时候追加 hello 到内容的前面
		if res.StatusCode != 200 {

			// 使用 ioutil.ReadAll方法获取返回的相应内容
			oldPayload, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			// 修改返回信息的内容，在响应内容的前面添加 hello
			newPayload := []byte("hello" + string(oldPayload))
			// 把修改之后的内容回写到 res.Body 中
			//Body io.ReadCloser
			// ReadCloser is the interface that groups the basic Read and Close methods.
			//type ReadCloser interface {
			//	Reader
			//	Closer
			//}
			//type Reader interface {
			//	Read(p []byte) (n int, err error)
			//}
			//type Closer interface {
			//	Close() error
			//}

			//type nopCloser struct {
			//	io.Reader
			//}

			// NopCloser returns a ReadCloser with a no-op Close method wrapping
			// the provided Reader r.
			//func NopCloser(r io.Reader) io.ReadCloser {
			//	return nopCloser{r}
			//}
			res.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			// 设置新的ContentLength
			res.ContentLength = int64(len(newPayload))
			// 往Header中写ContentLength，告诉客户端该读取多长
			res.Header.Set("Content-Length", fmt.Sprint(len(newPayload)))
		}
		return nil
	}

	// 把上面定义的处理逻辑函数传入ReverseProxy结构体并返回
	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyFunc,
	}
}

func singleJoiningSlash(a, b string) string {
	// 获取 a 的后缀
	aslash := strings.HasSuffix(a, "/")
	// 获取 b 的前缀
	bslash := strings.HasPrefix(b, "/")
	switch {
	// 如果两个都能取到 "/" ，就把 b 的去掉进行合并
	case aslash && bslash:
		return a + b[1:]
	// 如果两个都没有 "/"　，添加一个
	case !aslash && !bslash:
		return a + "/" + b
	}
	// 默认直接合并
	return a + b
}
