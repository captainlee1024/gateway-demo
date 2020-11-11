package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件

// HandlerFunc 回调函数
type HandlerFunc func(*SliceRouterContext)

// SliceRouter router 结构体
type SliceRouter struct {
	groups []*SliceGroup
}

// SliceGroup Group 结构体
type SliceGroup struct {
	*SliceRouter
	path     string
	handlers []HandlerFunc
}

// SliceRouterContext router 上下文
type SliceRouterContext struct {
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8 // 用来操作路由里面回调的数组
}

func newSliceRouterContext(rw http.ResponseWriter,
	req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}

	// 最长url前缀匹配
	matchURLLen := 0
	// 遍历路由器里的路由，挑出一个最匹配的
	for _, group := range r.groups {
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchURLLen {
				matchURLLen = pathLen
				*newSliceGroup = *group // 浅拷贝数组指针
			}
		}
	}

	// 利用最匹配的前缀，请求的上下文，responsewrite ， request来构建router上下文
	c := &SliceRouterContext{
		Rw:         rw,
		Req:        req,
		SliceGroup: newSliceGroup,
		Ctx:        req.Context(),
	}
	c.Reset()
	return c
}

// Group 创建 Group
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		SliceRouter: g,
		path:        path,
	}
}

// Use 构造回调方法
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	// 把中间件添加到路由组里面
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	// 判断该路由组是否添加到了路由器中，如果没有就添加到路由器的路由数组中
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}

// SliceRouterHandler xx
type SliceRouterHandler struct {
	coreFunc func(*SliceRouterContext) http.Handler
	router   *SliceRouter
}

// ServeHTTP xx
func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter,
	req *http.Request) {
	c := newSliceRouterContext(rw, req, w.router)
	if w.coreFunc != nil {
		c.handlers = append(c.handlers,
			func(c *SliceRouterContext) {
				w.coreFunc(c).ServeHTTP(rw, req)
			})
	}
	c.Reset()
	c.Next()

}

// NewSliceRouterHandler xx
func NewSliceRouterHandler(coreFunc func(*SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// NewSliceRouter 构造 router
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

// Get 获取上下文属性值
func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

// Set 设置上下文属性值
func (c *SliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

// Next 从最先加入中间件开始回调
func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 跳出中间件方法
func (c *SliceRouterContext) Abort() {
	c.index = abortIndex
}

// IsAborted 是否跳过了回调
func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

// Reset 重置回调
func (c *SliceRouterContext) Reset() {
	c.index = -1
}
