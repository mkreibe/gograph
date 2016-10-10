package graph

import (
	"encoding/json"
	"io/ioutil"
)

// jsonNode is the json representation for the node.
type jsonNode map[string]interface{}

// jsonEdge os the json connections.
type jsonEdge []interface{}

// jsonGraphRepresentation defines the graph representation.
type jsonGraphRepresentation struct {
	Type       Type              `json:"type"`
	Attributes map[string]string `json:"attributes"`
	Nodes      []jsonNode        `json:"nodes"`
	Edges      []jsonEdge        `json:"edges"`
}

// LoadFileGraph will load the graph from a file.
func LoadFileGraph(fileName string) (graph *Graph, err error) {

	var file []byte
	if file, err = ioutil.ReadFile(fileName); err == nil {
		jsonGraph := jsonGraphRepresentation{}
		if err = json.Unmarshal(file, &jsonGraph); err == nil {
			if graph, err = NewGraph(jsonGraph.Type); err == nil {

				// 1 - Load the attributes
				for key, value := range jsonGraph.Attributes {
					graph.Attributes.Set(key, value)
				}

				// 2 - Load the nodes
				for _, n := range jsonGraph.Nodes {

					nid := n["id"]

					if nid != nil {
						var node *Node
						if node, err = graph.AddNode(nid.(string)); err != nil {
							break
						}

						for attrName, attrVal := range n {
							if attrName != "id" {
								node.Attributes.Set(attrName, attrVal)
							}
						}
					}
				}

				// 3 - Load the connections
				if err == nil {
					for _, edge := range jsonGraph.Edges {

						var nodes []string
						attr := NewAttributeCollection()

						for _, item := range edge {
							switch item.(type) {
							case string:
								nodes = append(nodes, item.(string))
							case map[string]interface{}:
								attr.Merge(item.(map[string]interface{}), true)
							}
						}

						if _, err = graph.AddEdge(attr, nodes); err != nil {
							break
						}
					}
				}
			}
		}
	}

	return
}
