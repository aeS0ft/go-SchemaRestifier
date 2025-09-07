package datastructures

import "sync"

type Node struct {
	Name     string
	Fields   []*Fields
	Children []*Node
	Mu       sync.Mutex
	// updatable later
	Hidden bool
}

func IsNodeEmpty(n Node) bool {
	return n.Name == "" && len(n.Fields) == 0 && len(n.Children) == 0
}
