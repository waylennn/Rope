package registry

// Service service abstract
type Service struct {
	Name  string
	Nodes []*Nodes
}

// Nodes service node abstract
type Nodes struct {
	ID     string `json:"id"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}
