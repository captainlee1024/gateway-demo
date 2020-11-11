package tcpmiddleware

import (
	"context"
	"math"
	"net"

	"github.com/captainlee1024/gateway-demo/proxy/tcpproxy"
)

const abortIndex int8 = math.MaxInt8 / 2 // 最多63个中间件

// TCPHandlerFunc description
type TCPHandlerFunc func(*TCPSliceRouterContext)

// TCPSliceRouter router 结构体
type TCPSliceRouter struct {
	groups []*TCPSliceGroup
}

// TCPSliceGroup group 结构体
type TCPSliceGroup struct {
	*TCPSliceRouter
	path     string
	handlers []TCPHandlerFunc
}

// TCPSliceRouterContext 上下文
type TCPSliceRouterContext struct {
	conn net.Conn
	Ctx  context.Context
	*TCPSliceGroup
	index int8
}

func newTCPSliceRouterContext(ctx context.Context, conn net.Conn,
	r *TCPSliceRouter) *TCPSliceRouterContext {
	newTCPSliceGroup := &TCPSliceGroup{}
	// TODO 为什么要浅拷贝使用第一个分组
	*newTCPSliceGroup = *r.groups[0] // 浅拷贝数组指针，只会使用第一个分组
	c := &TCPSliceRouterContext{
		conn:          conn,
		TCPSliceGroup: newTCPSliceGroup,
		Ctx:           ctx,
	}
	c.Reset()
	return c
}

// Get 获取上下文的值
func (c *TCPSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

// Set 设置上下文的值
func (c *TCPSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

// TCPSliceRouterHandler description
type TCPSliceRouterHandler struct {
	coreFunc func(*TCPSliceRouterContext) tcpproxy.TCPHandler
	router   *TCPSliceRouter
}

// ServeTCP description
func (w *TCPSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTCPSliceRouterContext(ctx, conn, w.router)
	// TODO 每次调用都会进行append 那么这个方法切片会不会无限增大呢
	c.handlers = append(c.handlers, func(c *TCPSliceRouterContext) {
		w.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset() // 给索引设置成 -1 这样保证在请求的时候是从中间件的第一个开始调用的
	c.Next()
}

// NewTCPSliceRouterHandler 创建 TCP 路由处理器
func NewTCPSliceRouterHandler(coreFunc func(*TCPSliceRouterContext) tcpproxy.TCPHandler,
	router *TCPSliceRouter) *TCPSliceRouterHandler {
	return &TCPSliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// NewTCPSliceRouter 构造 router
func NewTCPSliceRouter() *TCPSliceRouter {
	return &TCPSliceRouter{}
}

// Group 创建 Group
func (g *TCPSliceRouter) Group(path string) *TCPSliceGroup {
	if path != "/" {
		panic("only accept path=/")
	}
	return &TCPSliceGroup{
		TCPSliceRouter: g,
		path:           path,
	}
}

// Use 构造回调方法
func (g *TCPSliceGroup) Use(middlewares ...TCPHandlerFunc) *TCPSliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.TCPSliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.TCPSliceRouter.groups = append(g.TCPSliceRouter.groups, g)
	}
	return g
}

// Next 从最先加入的中间件开始回调
func (c *TCPSliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 跳出中间件方法
func (c *TCPSliceRouterContext) Abort() {
	c.index = abortIndex
}

// IsAbort 是否跳过了回调方法
func (c *TCPSliceRouterContext) IsAbort() bool {
	return c.index >= abortIndex
}

// Reset 重置回调
func (c *TCPSliceRouterContext) Reset() {
	c.index = -1
}
