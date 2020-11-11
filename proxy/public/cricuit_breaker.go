package public

import (
	"log"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

// ConfCircuitBreaker 设置熔断
func ConfCircuitBreaker(openStream bool) {
	hystrix.ConfigureCommand("common", hystrix.CommandConfig{
		Timeout:                1000, // 单次请求超时时间
		MaxConcurrentRequests:  1,
		SleepWindow:            5000, // 熔断之后多久开始尝试服务是否可用
		RequestVolumeThreshold: 1,
		ErrorPercentThreshold:  1,
	})

	if openStream {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		go func() {
			err := http.ListenAndServe(net.JoinHostPort("", "2001"), hystrixStreamHandler)
			log.Fatal(err)
		}()
	}
}
