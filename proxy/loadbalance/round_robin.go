package loadbalance

import (
	"errors"
	"fmt"
	"strings"
)

// RoundRobinBalance 轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []string
	// 观察主题
	conf Conf
}

// Add 添加服务器
// 可以一次添加多个服务地址
// func (r *RoundRobinBalance) Add(params ...string) error {
// 	if len(params) == 0 {
// 		return errors.New("param len 1 at least")
// 	}
// 	for _, addr := range params {
// 		r.rss = append(r.rss, addr)
// 	}
// 	return nil
// }

// Add 添加服务
// 一次添加一个服务地址
func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

// Next 轮询方式获取服务器
func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = r.curIndex % lens
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

// Get 对Next方法进一步封装
func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

// SetConf 配置观察主题
func (r *RoundRobinBalance) SetConf(conf Conf) {
	r.conf = conf
}

// Update 更新，通知其他...
func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*ZKConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
	if conf, ok := r.conf.(*CheckConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
