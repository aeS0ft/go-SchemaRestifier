package datastructures

type Node struct {
	Name     string
	Fields   []*Fields
	Children []*Node
}

func IsNodeEmpty(n Node) bool {
	return n.Name == "" && len(n.Fields) == 0 && len(n.Children) == 0
}
