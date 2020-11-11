package main

import (
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func Test_main(t *testing.T) {
	// 启动一个流的服务器，来统计熔断降级限流的结果，统计数据会实时发送到上面
	// 可接入控制面板，可视化的呈现出来
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8074", hystrixStreamHandler)

	hystrix.ConfigureCommand("aaa", hystrix.CommandConfig{
		Timeout:               1000, // 单次请求超时时间 默认时间是1000毫秒
		MaxConcurrentRequests: 1,    // command的最大并发量 默认值是10，超过开启熔断
		SleepWindow:           5000, // 熔断后多久去尝试服务是否可用，由打开到半打开的时间间隔
		// 一个统计窗口10秒内请求数量。达到这个请求数量后才去判断是否要开启熔断。默认值是20
		RequestVolumeThreshold: 1,
		// 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
		// 默认值是50
		ErrorPercentThreshold: 1,
	})

	// 熔断器和业务逻辑整合
	for i := 0; i < 10000; i++ {
		// 异步调用 hystrix.Go

		// 同步调用hystrix.Do
		// aaa：熔断器的名称，必须跟上面的配置中的熔断器的名称相同
		// 后面两个参数一个是业务逻辑函数，一个是降级操作的函数，这里降级操作先不设置
		// 业务逻辑的方法内部出错就会执行后面的降级函数，如果降级出错，那么就输出错误到err里
		err := hystrix.Do("aaa", func() error {
			// test case 1 并发测试
			// 设置第一次请求总是返回错误
			if i == 0 {
				return errors.New("service error")
			}
			// test case 2 超时测试
			// time.Sleep(2 * time.Second)
			log.Println("do services")
			return nil
		}, nil)
		if err != nil {
			log.Println("hystrix err:" + err.Error())
			time.Sleep(1 * time.Second)
			log.Println("sleep 1 second")
		}
	}
	time.Sleep(100 * time.Second)
}
