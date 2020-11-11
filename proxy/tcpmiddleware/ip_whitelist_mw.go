package tcpmiddleware

import (
	"strings"
)

// IPWhiteListMiddleWare TCP 权限认证中间件
// 定义了一份IP白名单，这里是设置127.0.0.1这个 IP　可以通过验证
func IPWhiteListMiddleWare() func(c *TCPSliceRouterContext) {
	return func(c *TCPSliceRouterContext) {
		// 从 conn 中获取来源 IP
		remoteAddr := c.conn.RemoteAddr().String()
		// 判断来源 IP 是不是 白名单（这里是127.0.0.1）里的
		// 如果是，就放行，代用下一个中间件方法，如果不行就跳出调用，写入相应提示信息返回
		// remoteAddr 是有端口的，所以这里比较前缀
		if strings.HasPrefix(remoteAddr, "127.0.0.1") {
			// 测试白名单不通过的情况
			// if strings.HasPrefix(remoteAddr, "127.0.0.2") {
			c.Next()
		} else {
			c.Abort()
			c.conn.Write([]byte("ip whitelist auth invalid"))
			c.conn.Close()
			// fmt.Println("auth failed" + c.conn.RemoteAddr().String()) // 测试的时候提供观察
		}
	}
}
