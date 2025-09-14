package datastructures

import "sync"

type Node struct {
	Name     string
	Fields   []*Field
	Children []*Node
	Mu       sync.Mutex
	// updatable later
	Hidden bool
}

func IsNodeEmpty(n *Node) bool {
	return n.Name == "" && len(n.Fields) == 0 && len(n.Children) == 0
}
func IsNodeLeaf(n *Node) bool {
	return len(n.Fields) == 0 && len(n.Children) == 0
}
func AllLeafsExhausted(node *Node) bool {
	if node == nil {
		return true
	}
	if len(node.Children) == 0 {
		// Leaf node reached
		return true
	}
	for _, child := range node.Children {
		if !AllLeafsExhausted(child) {
			return false
		}
	}
	return true
}
