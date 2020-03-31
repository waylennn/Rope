package loadbalance

import (
	"context"
	"cope/registry"
)

//LoadBalance 负载均衡
type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Nodes) (node *registry.Nodes, err error)
}

const (
	//DefaultNodeWeight 默认权重
	DefaultNodeWeight = 100
)
