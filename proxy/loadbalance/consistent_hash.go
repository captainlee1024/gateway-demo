package loadbalance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Hash 哈希函数
type Hash func(data []byte) uint32

// UInt32Slice 存储排序节点
type UInt32Slice []uint32

// Len 返回长度
func (s UInt32Slice) Len() int {
	return len(s)
}

// Less 比较函数
func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

// Swap 交换函数
func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ConsistentHashBanlance 一致性Hash负载均衡
type ConsistentHashBanlance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int               // 复制因子
	keys     UInt32Slice       // 已排序的节点 hash 切片
	hashMap  map[uint32]string // 节点哈希和 Key 的 map，键是 hash 值，值是节点 key

	// 观察主题
	conf Conf
}

// NewConsistentHashBanlance 创建
func NewConsistentHashBanlance(replicas int, fn Hash) *ConsistentHashBanlance {
	m := &ConsistentHashBanlance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		// 最多 32 位，保证是一个 2^32-1 环
		// 默认使用循环冗余的 crc32 方法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty 验证是否为空
func (c *ConsistentHashBanlance) IsEmpty() bool {
	return len(c.keys) == 0
}

// Add 方法用来添加缓存节点，参数为节点 key，比如使用 IP
func (c *ConsistentHashBanlance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	// 加锁
	c.mux.Lock()
	defer c.mux.Unlock()
	// 结合复制因子计算所有虚拟节点的 hash 值，并存入 m.keys 中
	//同时在 m.hashMap 中保存哈希值和 key 的映射
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}

	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(c.keys)
	return nil
}

// Get 方法根据给定的对象获取最靠近它的那个节点
func (c *ConsistentHashBanlance) Get(key string) (string, error) {
	if c.IsEmpty() {
		fmt.Println("consistent get ...")
		return "", errors.New("Node is empty")
	}
	hash := c.hash([]byte(key))
	// 通过二分查找获取最优节点，第一个“服务器hash”值大于“数据hash”值的就是最优“服务器节点”
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	// 如果查找结果大于服务器节点hash数组的最大索引，表示此时此刻对象hash值位于最后一个节点
	// 那么放入第一个节点中
	if idx == len(c.keys) {
		idx = 0
	}
	// 读取数据　添加读取的锁
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

// SetConf 设置
func (c *ConsistentHashBanlance) SetConf(conf Conf) {
	c.conf = conf
}

// Update 更新配置，通知其他...
func (c *ConsistentHashBanlance) Update() {
	if conf, ok := c.conf.(*ZKConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		c.keys = nil
		// c.hashMap = nil
		c.hashMap = map[uint32]string{}
		for _, ip := range conf.GetConf() {
			c.Add(strings.Split(ip, ",")...)
		}
	}
	if conf, ok := c.conf.(*CheckConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		c.keys = nil
		c.hashMap = map[uint32]string{}
		for _, ip := range conf.GetConf() {
			c.Add(strings.Split(ip, "")...)
		}
	}

}
