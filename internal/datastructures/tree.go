package datastructures

type Node struct {
	Name     string
	Fields   []*Fields
	Children []*Node
}
