package proxy

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/captainlee1024/gateway-demo/proxy/loadbalance"
	"github.com/captainlee1024/gateway-demo/proxy/tcpmiddleware"
)

// NewTCPLoadBalanceReverseProxy 创建TCP代理负载均衡器
func NewTCPLoadBalanceReverseProxy(c *tcpmiddleware.TCPSliceRouterContext,
	lb loadbalance.LoadBalance) *TCPReverseProxy {
	return func() *TCPReverseProxy {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		return &TCPReverseProxy{
			ctx:             c.Ctx,
			Addr:            nextAddr,
			KeepAlivePeriod: time.Second,
			DialTimeout:     time.Second,
		}
	}()
}

var defaultDialer = new(net.Dialer)

// TCPReverseProxy TCP 反向代理配置
type TCPReverseProxy struct {
	ctx             context.Context // 请求上下文 单次请求单独设置
	Addr            string
	KeepAlivePeriod time.Duration // 长连接设置
	DialTimeout     time.Duration // 拨号超时时间
	// 拨号的上下文
	DialContext func(ctx context.Context, network, address string) (net.Conn, error)
	// 拨号的错误
	OnDialError          func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
}

func (dp *TCPReverseProxy) dialTimeout() time.Duration {
	if dp.DialTimeout > 0 {
		return dp.DialTimeout
	}
	return 10 * time.Second
}

func (dp *TCPReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if dp.DialContext != nil {
		return dp.DialContext
	}
	return (&net.Dialer{
		Timeout:   dp.DialTimeout,     // 连接超时
		KeepAlive: dp.KeepAlivePeriod, // 设置连接的检测时长
	}).DialContext
}

func (dp *TCPReverseProxy) keepAlivePeriod() time.Duration {
	if dp.KeepAlivePeriod != 0 {
		return dp.KeepAlivePeriod
	}
	return time.Minute
}

// ServeTCP 实现 TCPHandler 接口，实现 TCP 反向代理逻辑
func (dp *TCPReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	// 设置连接超时
	var cancel context.CancelFunc
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, dp.DialTimeout)
	}
	// 设置拨号上下文 类型 地址
	dst, err := dp.dialContext()(ctx, "tcp", dp.Addr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		dp.onDialError()(src, err)
		return
	}

	defer func() { go dst.Close() }() // 记得关闭下游连接

	// 设置 dst （下游目标地址）的 KeepAlive 时间，在数据写入之前
	if ka := dp.keepAlivePeriod(); ka > 0 {
		if c, ok := dst.(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
	}

	errc := make(chan error, 1)
	go dp.proxyCopy(errc, src, dst)
	go dp.proxyCopy(errc, dst, src)
	<-errc
}

func (dp *TCPReverseProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if dp.OnDialError != nil {
		return dp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		log.Printf("tcpproxy: for incoming conn %v, error dialing %q: %v\n",
			src.RemoteAddr(), dp.Addr, dstDialErr)
		src.Close()
	}
}

func (dp *TCPReverseProxy) proxyCopy(errc chan<- error, dst, src net.Conn) {
	_, err := io.Copy(dst, src) // 从src读取并写入dst
	errc <- err
}
