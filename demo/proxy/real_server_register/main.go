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

	"github.com/captainlee1024/gateway-demo/proxy/zookeeper"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()
	time.Sleep(time.Second * 2)
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()
	time.Sleep(time.Second * 2)

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// RealServer 服务
type RealServer struct {
	Addr string
}

// Run 启动服务
func (r *RealServer) Run() {
	log.Println("Starting httpserver at" + r.Addr)
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则和处理器
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)

	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		// 连接 zk 服务器
		zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Printf("connect error: %s\n", err)
		}
		defer zkManager.Close()

		// 注册 zk 临时节点
		// 把当前服务的地址加端口注册到 /real_server
		err = zkManager.RegistServerPath("/real_server", r.Addr)
		if err != nil {
			fmt.Printf("register error: %s\n", err)
		}

		// 获取注册中心列表
		zlist, err := zkManager.GetServerListByPath("/real_server")
		fmt.Println(zlist)
		log.Fatal(server.ListenAndServe())
	}()

}

// HelloHandler Hello 处理器
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
