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
)

func main() {
	rs1 := &RealServer{
		Addr: "127.0.0.1:2003",
	}
	rs1.Run()
	rs2 := &RealServer{
		Addr: "127.0.0.1:2004",
	}
	rs2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}

// RealServer 请求的真实服务器的结构体。
type RealServer struct {
	Addr string
}

// Run 启动
func (r *RealServer) Run() {
	log.Println("Start httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_strip_uri/test_strip_uri/aa", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// HelloHandler hello处理器
func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	// 127.0.0.1:8080/abc?page=1
	// r.Addr = 127.0.0.1:8080
	// req.URL.Path = /abc
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)

	realIP := fmt.Sprintf("RemoteAddr=%s, X-Forwarded-For=%v, X-Real-IP=%v\n",
		req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("headers := %v\n", req.Header)

	io.WriteString(w, upath)
	w.Write([]byte(realIP))
	w.Write([]byte(header))
	w.Write([]byte("hello\n"))

}

// ErrorHandler error处理器
func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	//upath := "error handler"
	w.WriteHeader(500)
	w.Write([]byte("error handler\n"))
}

// TimeOutHandler 超时测试
func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "Timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
	//w.Write([]byte("timeout handler\n"))
}
