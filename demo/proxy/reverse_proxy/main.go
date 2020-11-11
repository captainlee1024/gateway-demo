package main

/*
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
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

	// rs2 := "http://127.0.0.1:2004"
	// url2, err2 := url.Parse(rs2)
	// if err2 != nil {
	// 	log.Println(err2)
	// }

	// 3. 调用自己的 NewSingleHostReverseProxy
	proxy := NewSingleHostReverseProxy(url1)
	//urls := []*url.URL{url1, url2}
	//proxy := NewSingleHostReverseProxy(urls)

	log.Println("starting httpserver at " + addr)
	// 4. 开启代理服务器，并代理目标服务器
	// 把proxy当做路由处理器来使用，因为ReverseProxy实现了ServeHTTP方法
	log.Fatal(http.ListenAndServe(addr, proxy))
}

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   time.Second * 30, // 连接超时
		KeepAlive: time.Second * 30, // 长连接时间
	}).DialContext,
	MaxIdleConns:          100,              // 最大空闲连接时间
	IdleConnTimeout:       time.Second * 90, // 空闲超时时间
	TLSHandshakeTimeout:   time.Second * 10, // tls握手超时时间
	ExpectContinueTimeout: time.Second * 1,  // 100-continue 超时时间
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
		// url地址重写： 重写前 /dir 重写后： /base/dir
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
		// 只在第一代里中设置此 header 头
		//req.Header.Set("X-Real-Ip", req.RemoteAddr)
	}

	// 更改内容
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
*/

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var addr = "127.0.0.1:2002"

func main() {
	rs1 := "http://www.baidu.com"
	// rs1 := "http://127.0.0.1:2003"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}

	rs2 := "http://www.baidu.com"
	// rs2 := "http://127.0.0.1:2004"
	url2, err2 := url.Parse(rs2)
	if err2 != nil {
		log.Println(err2)
	}
	urls := []*url.URL{url1, url2}
	proxy := NewMultipleHostsReverseProxy(urls)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, //连接超时
		KeepAlive: 30 * time.Second, //长连接超时时间
	}).DialContext,
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue 超时时间
}

// NewMultipleHostsReverseProxy 基于源码实现修改相应内容
func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		//url_rewrite
		//127.0.0.1:2002/dir/abc ==> 127.0.0.1:2003/base/abc ??
		//127.0.0.1:2002/dir/abc ==> 127.0.0.1:2002/abc
		//127.0.0.1:2002/abc ==> 127.0.0.1:2003/base/abc
		re, _ := regexp.Compile("^/dir(.*)")
		req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")

		//随机负载均衡
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//todo 部分章节补充1
		//todo 当对域名(非内网)反向代理时需要设置此项。当作后端反向代理时不需要
		req.Host = target.Host

		// url地址重写：重写前：/aa 重写后：/base/aa
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
		//只在第一代理中设置此header头
		//req.Header.Set("X-Real-Ip", req.RemoteAddr)
	}
	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//请求以下命令：curl 'http://127.0.0.1:2002/error'
		//todo 部分章节功能补充2
		//todo 兼容websocket
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		var payload []byte
		var readErr error

		//todo 部分章节功能补充3
		//todo 兼容gzip压缩
		if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			payload, readErr = ioutil.ReadAll(gr)
			resp.Header.Del("Content-Encoding")
		} else {
			payload, readErr = ioutil.ReadAll(resp.Body)
		}
		if readErr != nil {
			return readErr
		}

		//异常请求时设置StatusCode
		if resp.StatusCode != 200 {
			payload = []byte("StatusCode error:" + string(payload))
		}

		//todo 部分章节功能补充4
		//todo 因为预读了数据所以内容重新回写
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}
	//错误回调 ：关闭real_server时测试，错误回调
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}
	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc}
}

func singleJoiningSlash(a, b string) string {
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
