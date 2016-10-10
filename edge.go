package graph

import "errors"

// Edge is the connections between nodes.
type Edge struct {

	// Links to the nodes
	ID string

	// Attributes asssociated with this edge.
	Attributes AttributeCollection

	// Nodes this edge is connected to. NOTE: This could be used for hyper-graphs.
	Nodes map[string]*Node
}

// NewEdge creates a new edge.
func NewEdge(attrs AttributeCollection, connects []*Node) (edge *Edge, err error) {

	edge = &Edge{}
	edge.Attributes = attrs
	edge.Nodes = make(map[string]*Node)

	// make the id.
	if edge.Attributes.Contains("id") {
		id, _ := edge.Attributes.Get("id")
		if idStr := id.(string); len(idStr) > 0 {
			edge.ID = idStr
			edge.Attributes.Remove("id")
		} else {
			err = errors.New("Invalid edge id")
		}
	} else {
		for _, node := range connects {
			if len(edge.ID) == 0 {
				edge.ID = node.ID
			} else {
				edge.ID += "-" + node.ID
			}
		}
	}

	if len(connects) < 2 {
		err = errors.New("Edges must have atleast two nodes to connect")
	}

	// add the edge to the nodes collection.
	if err == nil {
		for _, node := range connects {
			edge.Nodes[node.ID] = node
			node.Edges[edge.ID] = edge
		}
	} else {
		edge = nil
	}

	return
}
