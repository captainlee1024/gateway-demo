package tcpmiddleware

import (
	"fmt"

	"github.com/captainlee1024/gateway-demo/proxy/public"
)

// FlowCountMiddleWare TCP 限流中间件
func FlowCountMiddleWare(counter *public.FlowCountService) func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		counter.Increase()
		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)
		c.Next()
	}
}
