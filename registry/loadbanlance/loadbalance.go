package loadbanlance

import (
	"context"
	"cope/registry"
	"errors"
)

var (
	//ErrNotHaveNodes ..
	ErrNotHaveNodes = errors.New("not have node")
)

const (
	//DefaultNodeWeight set defalut node weight
	DefaultNodeWeight = 100
)

// LoadBalance ...
type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Nodes) (node *registry.Nodes, err error)
}
