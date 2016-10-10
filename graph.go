package graph

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

const (

	// GraphDirected is a directed graph.
	GraphDirected Type = "directed"

	// GraphUndirected is an undirected graph.
	GraphUndirected Type = "undirected"
)

// Type describes the known types of graphs.
type Type string

// Graph is the collection of nodes and edges.
type Graph struct {

	// Attributes for this graph.
	Attributes AttributeCollection

	// Type of graph this is.
	Type Type

	// Edges contained within this graph.
	Edges map[string]*Edge

	// Nodes contained within this graph.
	Nodes map[string]*Node
}

// NewGraph will create a new graph.
func NewGraph(graphType Type) (graph *Graph, err error) {
	graph = &Graph{
		Type:       graphType,
		Attributes: NewAttributeCollection(),
		Edges:      make(map[string]*Edge),
		Nodes:      make(map[string]*Node),
	}
	return
}

// AddNode to the graph.
func (graph *Graph) AddNode(id string) (node *Node, err error) {
	if node, err = NewNode(id); err == nil {
		graph.Nodes[id] = node
	}

	return
}

// AddEdge will create and attach the edge to the appropriate nodes.
func (graph *Graph) AddEdge(attrs AttributeCollection, nodeIds []string) (edge *Edge, err error) {
	var nodes []*Node

	for _, id := range nodeIds {
		nodes = append(nodes, graph.Nodes[id])
	}

	if edge, err = NewEdge(attrs, nodes); err == nil {
		graph.Edges[edge.ID] = edge
	}

	return
}

// HasConnection will check if the connection exists.
func (graph *Graph) HasConnection(source string, target string) (result bool, err error) {

	switch graph.Type {
	case GraphUndirected:
		{
			for _, edge := range graph.Edges {

				hasSource := false
				hasTarget := false

				for _, node := range edge.Nodes {

					if node.ID == source {
						hasSource = true
					}

					if node.ID == target {
						hasTarget = true
					}

					if hasSource && hasTarget {
						result = true
						break
					}
				}

				if result {
					break
				}
			}
			break
		}
	default:
		err = fmt.Errorf("Unknown graph type: %s", graph.Type)
	}
	return
}

// IterNodes will iterate over all the nodes in order of the attribute name.
func (graph *Graph) IterNodes(attr string, iterFunc func(node *Node) (cont bool, err error)) (err error) {

	if iterFunc == nil {
		err = errors.New("No iteration function")
	} else {
		nodeMap := map[interface{}][]*Node{}
		for _, node := range graph.Nodes {
			if a, ok := node.Attributes.Get(attr); ok {
				nodeMap[a] = append(nodeMap[a], node)
			}
		}

		var attrsValues sort.Interface
		for key := range nodeMap {
			switch key.(type) {
			case string:
				if attrsValues == nil {
					attrsValues = sort.StringSlice{}
				}
				attrsValues = append(attrsValues.(sort.StringSlice), key.(string))
			case int:
				if attrsValues == nil {
					attrsValues = sort.IntSlice{}
				}
				attrsValues = append(attrsValues.(sort.IntSlice), key.(int))
			case float64:
				if attrsValues == nil {
					attrsValues = sort.Float64Slice{}
				}
				attrsValues = append(attrsValues.(sort.Float64Slice), key.(float64))
			default:
				err = errors.New("Non sortable type. Requires string, int or float64")
				break
			}
		}

		if err == nil {
			sort.Sort(attrsValues)

			attrs := reflect.ValueOf(attrsValues)
			for i := 0; i < attrs.Len(); i++ {
				value := attrs.Index(i).Interface()
				cont := false
				for _, node := range nodeMap[value] {
					if cont, err = iterFunc(node); err != nil || !cont {
						break
					}
				}

				if err != nil || !cont {
					break
				}
			}
		}
	}
	return
}
