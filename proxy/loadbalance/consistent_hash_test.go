package loadbalance

import (
	"fmt"
	"testing"
)

func TestConsistentHashBanlance(t *testing.T) {
	// 虚拟节点个数可以根据业务实际情况去设置
	rb := NewConsistentHashBanlance(10, nil)
	rb.Add("127.0.0.1:2003")
	rb.Add("127.0.0.1:2004")
	rb.Add("127.0.0.1:2005")
	rb.Add("127.0.0.1:2006")
	rb.Add("127.0.0.1:2007")

	// url hash
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/error"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/changepwd"))

	// ip hash
	fmt.Println(rb.Get("127.0.0.1"))
	fmt.Println(rb.Get("192.268.0.1"))
	fmt.Println(rb.Get("127.0.0.1"))

}
