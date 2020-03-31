package loadbalance

import (
	"context"
	"cope/meta"
	"cope/registry"
	"math/rand"
	"reflect"

	"cope/errno.go"
)

// RandomBalance 加权随机负载均衡算法
type RandomBalance struct {
	name string
}

//NewRandomBalance ...
func NewRandomBalance() (bandombalance *RandomBalance) {
	return &RandomBalance{
		name: "random",
	}
}

//Name ...
func (r *RandomBalance) Name() (name string) {
	name = r.name
	return
}

//Select ...
func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Nodes) (node *registry.Nodes, err error) {
	if len(nodes) == 0 {
		err = errno.NotHaveInstance
		return
	}

	var totalWeight int
	for _, node := range nodes {
		if node.Weight == 0 {
			node.Weight = DefaultNodeWeight
		}
		totalWeight += node.Weight
	}

	curWeight := rand.Intn(totalWeight)
	index := 0

	for ind, node := range nodes {
		curWeight -= node.Weight
		if curWeight < 0 {
			index = ind
			break
		}
	}
	node = nodes[index]
	r.RemoveNode(ctx, node)
	return
}

//RemoveNode 移除已经使用过的元素
func (r *RandomBalance) RemoveNode(ctx context.Context, node *registry.Nodes) {
	rpcMeta := meta.GetRpcMeta(ctx)
	for k, v := range rpcMeta.AllNodes {
		if reflect.DeepEqual(node, v) {
			rpcMeta.AllNodes = append(rpcMeta.AllNodes[:k], rpcMeta.AllNodes[k+1:]...)
			break
		}
	}
}
