package loadbanlance

import (
	"context"
	"cope/registry"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	var nodes []*registry.Nodes
	temp := []int{50, 100, 150}
	for i := 0; i < 3; i++ {
		node := &registry.Nodes{
			IP:     fmt.Sprintf("192.0.0.%v", i),
			Port:   9000,
			Weight: temp[i%3],
		}
		nodes = append(nodes, node)
	}

	newBandomBalance := NewRandomBalance()

	tempMap := make(map[string]int)

	for i := 0; i < 1000; i++ {
		node, err := newBandomBalance.Select(context.TODO(), nodes)
		if err != nil {
			panic(err)
		}

		tempMap[fmt.Sprintf("ip:%v;weight:%v", node.IP, node.Weight)]++

	}

	fmt.Println(tempMap)
}
