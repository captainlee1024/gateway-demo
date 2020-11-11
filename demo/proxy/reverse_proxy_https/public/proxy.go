package public

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/captainlee1024/gateway-demo/demo/proxy/reverse_proxy_https/testdata"
	"golang.org/x/net/http2"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   time.Second * 30, // 连接超时
		KeepAlive: time.Second * 30, // 长连接超时时长
	}).DialContext,
	// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	TLSClientConfig: func() *tls.Config {
		pool := x509.NewCertPool()
		caCertPath := testdata.Path("ca.crt")
		caCrt, _ := ioutil.ReadFile(caCertPath)
		pool.AppendCertsFromPEM(caCrt)
		return &tls.Config{RootCAs: pool}
	}(),
	MaxIdleConns:          100,              // 最大空闲连接
	IdleConnTimeout:       time.Second * 90, // 空闲超时时间
	TLSHandshakeTimeout:   time.Second * 10, // tls 握手超时时间
	ExpectContinueTimeout: time.Second * 1,  // 100-continue 超时时间
}

// NewMultipleHostsReverseProxy 支持https代理
func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery
		fmt.Println("target.Scheme")
		fmt.Println(target.Scheme)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}
	http2.ConfigureTransport(transport)
	return &httputil.ReverseProxy{Director: director, Transport: transport}
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
