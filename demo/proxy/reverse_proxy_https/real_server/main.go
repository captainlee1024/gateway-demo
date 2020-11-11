package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/captainlee1024/gateway-demo/demo/proxy/reverse_proxy_https/testdata"
	"golang.org/x/net/http2"
)

/*
证书签名生成方式:

//CA私钥
openssl genrsa -out ca.key 2048
//CA数据证书
openssl req -x509 -new -nodes -key ca.key -subj "/CN=example1.com" -days 5000 -out ca.crt

//服务器私钥（默认由CA签发）
openssl genrsa -out server.key 2048
//服务器证书签名请求：Certificate Sign Request，简称csr（example1.com代表你的域名）
openssl req -new -key server.key -subj "/CN=example1.com" -out server.csr
//上面2个文件生成服务器证书（days代表有效期）
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000
*/

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:3003"}
	rs1.Run()
	rs2 := &RealServer{Addr: "127.0.0.1:3004"}
	rs2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}

// RealServer 结构体
type RealServer struct {
	Addr string
}

// Run 封装的启动服务器方法
func (r *RealServer) Run() {
	log.Println("Starting httpserver at" + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		// 把 http 服务器升级成 http2
		http2.ConfigureServer(server, &http2.Server{})
		server.ListenAndServeTLS(testdata.Path("server.crt"), testdata.Path("server.key"))
	}()
}

// HelloHandler 处理器
func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	io.WriteString(w, upath)
}

// ErrorHandler 错误处理器
func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}
