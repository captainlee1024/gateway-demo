package loadbalance

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// RandomBalance 存放服务器列表和当前使用服务器下标
type RandomBalance struct {
	curIndex int
	rss      []string
	// 观察主题
	conf Conf
}

// Add 添加服务器
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}

	for _, addr := range params {
		r.rss = append(r.rss, addr)
	}
	// addr := params[0]
	// r.rss = append(r.rss, addr)
	return nil
}

// Next 从服务器列表随机获取一台服务器
func (r *RandomBalance) Next() string {
	// 验证服务器数组是否为空
	if len(r.rss) == 0 {
		return ""
	}
	// 随机返回服务器
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

// Get 对 Next 进一步包装
func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

// SetConf 初始化配置
func (r *RandomBalance) SetConf(conf Conf) {
	r.conf = conf
}

// Update 更新配置
func (r *RandomBalance) Update() {
	if conf, ok := r.conf.(*ZKConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*CheckConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		// r.rss = nil
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}

	}
}
