package loadbanlance

import (
	"context"
	"cope/registry"
	"math/rand"
)

// RandomBalance ...
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

	var totalWeight int
	for _, node := range nodes {
		if node.Weight == 0 {
<<<<<<< HEAD
			node.Weight = DefaultNodeWeight
=======
			totalWeight = DefaultNodeWeight
>>>>>>> 4f008cdac20aeb2e116d997fc89c3568a6dca67a
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

	return
}
