package loadbalance

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	rb := &RoundRobinBalance{}
	// rb.Add(
	// 	"127.0.0.1:2003",
	// 	"127.0.0.1:2004",
	// 	"127.0.0.1:2005",
	// 	"127.0.0.1:2006",
	// 	"127.0.0.1:2007",
	// )
	rb.Add("127.0.0.1:2003")
	rb.Add("127.0.0.1:2004")
	rb.Add("127.0.0.1:2005")
	rb.Add("127.0.0.1:2006")
	rb.Add("127.0.0.1:2007")

	for i := 0; i < 10; i++ {
		fmt.Println(rb.Next())
	}
}
