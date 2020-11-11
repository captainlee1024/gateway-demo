package loadbalance

import (
	"fmt"
	"testing"
)

func TestLWeightRoundRObinBalance(t *testing.T) {
	br := &WeightRoundRobinBalance{}
	br.Add("127.0.0.1:2003", "4")
	br.Add("127.0.0.1:2004", "3")
	br.Add("127.0.0.1:2005", "2")

	for i := 0; i < 18; i++ {
		fmt.Println(br.Next())
	}
}
