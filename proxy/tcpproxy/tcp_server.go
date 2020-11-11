package tcpproxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ErrServerClosed
var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

// TCPHandler TCP 处理器
type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

// TCPServer xx
type TCPServer struct {
	Addr    string
	Handler TCPHandler
	err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

func (s *TCPServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

// Close 关闭连接
func (s *TCPServer) Close() error {
	atomic.StoreInt32(&s.inShutdown, 1)
	close(s.doneChan) // 关闭channel
	s.Close()         // 执行 listener 关闭
	return nil
}

// ListenAndServe 监听开启服务
func (s *TCPServer) ListenAndServe() error {
	if s.shuttingDown() {
		return ErrServerClosed
	}
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	addr := s.Addr
	if addr == "" {
		return errors.New("need addr")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(tcpKeepAliveListener{
		ln.(*net.TCPListener),
	})
}

// Serve 创建
func (s *TCPServer) Serve(l net.Listener) error {
	s.l = &onceCloseListener{Listener: l}
	defer s.l.Close() // 执行 listener 关闭
	if s.BaseCtx == nil {
		s.BaseCtx = context.Background()
	}
	baseCtx := s.BaseCtx
	ctx := context.WithValue(baseCtx, ServerContextKey, s)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-s.getDoneChan():
				return ErrServerClosed
			default:
			}
			fmt.Printf("accept fail., err:%v\n", e)
			continue
		}
		c := s.newConn(rw)
		go c.serve(ctx)
	}
}

func (s *TCPServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: s,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

func (s *TCPServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}
