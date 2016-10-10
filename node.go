package graph

import "errors"

// Node is the vertex of the graph.
type Node struct {

	// ID for this node.
	ID string

	// Attributes associated with this node.
	Attributes AttributeCollection

	// Edges this node is connected to.
	Edges map[string]*Edge
}

// NewNode creates a new node.
func NewNode(id string) (node *Node, err error) {

	if len(id) != 0 {
		node = &Node{
			ID:         id,
			Attributes: NewAttributeCollection(),
			Edges:      make(map[string]*Edge),
		}
	} else {
		err = errors.New("Node must have a valid id")
	}

	return
}

// Adj will return the adjacent nodes.
func (node *Node) Adj() (adjacent []*Node) {

	adjacent = []*Node{}

	for _, edge := range node.Edges {
		for _, n := range edge.Nodes {
			if n != node {
				adjacent = append(adjacent, n)
			}
		}
	}

	return
}
