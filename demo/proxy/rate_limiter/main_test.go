package main

import (
	"context"
	"log"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func Test_RateLimiter(t *testing.T) {
	l := rate.NewLimiter(1, 5)
	log.Println(l.Limit())
	for i := 0; i < 100; i++ {
		// 阻塞等待，直到取到一个 token
		log.Println("before wait")
		c, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
		defer cancelFunc()
		if err := l.Wait(c); err != nil {
			log.Println("limiter wait err:", err.Error())
		}
		log.Println("after wait")

		// 返回需要等待多久才有新的token，这样就可以等待指定时间执行任务
		r := l.Reserve()
		log.Println("reserve Delay:", r.Delay())

		// 判断当前是否可以取到token
		a := l.Allow()
		log.Println("allow:", a)
	}
}
