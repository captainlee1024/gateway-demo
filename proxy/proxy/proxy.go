package proxy

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/middleware"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   time.Second * 30, // 连接超时时间
		KeepAlive: time.Second * 30, // 长连接超时时间
	}).DialContext,
	MaxIdleConns:          100,              // 最大空闲连接
	IdleConnTimeout:       time.Second * 90, // 空闲超时时间
	TLSHandshakeTimeout:   time.Second * 10, // tls 握手超时时间
	ExpectContinueTimeout: time.Second * 1,  // 100-continue 状态码超时时间
}

// NewMultipleHostsReverseProxy 反向代理负载均衡
// func NewMultipleHostsReverseProxy(lb loadbalance.LoadBalance) *httputil.ReverseProxy {
func NewMultipleHostsReverseProxy(c *middleware.SliceRouterContext, targets []*url.URL) *httputil.ReverseProxy {

	// 请求协调者
	director := func(req *http.Request) {
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery
		// nextAddr, err := lb.Get(req.RemoteAddr)
		// // nextAddr, err := lb.Get(req.URL.String())
		// if err != nil {
		// 	log.Fatal("get next addr fail")
		// }
		// target, err := url.Parse(nextAddr)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoinSlash(target.Path, req.URL.Path)
		// TODO 当对域名（非内网）反向代理时需要设置此项，当作后端反向代理时不需要
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
		// if resp.StatusCode != 200 {
		// 	// 获取内容
		// 	oldPayload, err := ioutil.ReadAll(resp.Body)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// 追加内容
		// 	newPayload := []byte("StatusCode error:" + string(oldPayload))
		// 	resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
		// 	resp.ContentLength = int64(len(newPayload))
		// 	resp.Header.Set("Content", strconv.FormatInt(int64(len(newPayload)), 10))
		// }

		// TODO 兼容websocket
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		var payload []byte
		var readErr error

		// TODO 兼容 gzip 压缩
		if strings.Contains(resp.Header.Get("Connection-Encoding"), "gzip") {
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

		// 异常请求设置 StatueCode
		if resp.StatusCode != 200 {
			payload = []byte("StatusCode error:" + string(payload))
		}

		// TODO 因为预读了数据所以内容重新回写
		c.Set("status_code", resp.StatusCode)
		c.Set("payload", payload)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	// 错误回调：关闭real_server时测试，错误回调
	// 范围：transport.RoundTrip 发生错误，以及 ModifyResponse 发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		// TODO 如果是权重的负载则调整临时权重
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
