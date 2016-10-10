package graph

import (
	"errors"
	"fmt"
	"sort"
	"testing"
)

func AssertNodeAttributeValue(t *testing.T, graph *Graph, node string, attr string, expected interface{}) {
	test := fmt.Sprintf("Depth[%s]", node)

	if n, ok := graph.Nodes[node]; ok {
		if actual, ok := n.Attributes.Get(attr); ok {
			AssertT(t, test, expected, actual)
		} else {
			t.Error(test, " Invalid attribute.")
		}
	} else {
		t.Error(test, " Invalid node.")
	}
}

func TestBFSAlgorithm(t *testing.T) {
	t.Log("Graph loading.")
	if graph, err := LoadFileGraph("./data/ud_8n_9e.json"); err != nil {
		t.Error(err)
	} else {
		graph.BFSd("s")

		AssertNodeAttributeValue(t, graph, "r", "d", 1)
		AssertNodeAttributeValue(t, graph, "s", "d", 0)
		AssertNodeAttributeValue(t, graph, "t", "d", 2)
		AssertNodeAttributeValue(t, graph, "u", "d", 3)
		AssertNodeAttributeValue(t, graph, "v", "d", 2)
		AssertNodeAttributeValue(t, graph, "w", "d", 1)
		AssertNodeAttributeValue(t, graph, "x", "d", 2)
		AssertNodeAttributeValue(t, graph, "y", "d", 3)

		depths := map[int]map[string]int{
			0: map[string]int{"s": -1},
			1: map[string]int{"r": -1, "w": -1},
			2: map[string]int{"t": -1, "v": -1, "x": -1},
			3: map[string]int{"u": -1, "y": -1},
		}

		index := 0
		graph.IterNodes("d", func(node *Node) (cont bool, err error) {
			if i, ok := node.Attributes.Get("d"); ok {
				depths[i.(int)][node.ID] = index
				index++
				return true, nil
			}

			return false, errors.New("No 'd' attribute.")
		})

		next := 0
		for i := 0; i < len(depths); i++ {
			current := depths[i]

			var a []int
			for _, v := range current {
				a = append(a, v)
			}

			sort.Ints(a)
			for _, v := range a {
				if next != v {
					t.Errorf("Expected %d, got %d", next, v)
				}
				next++
			}
		}
	}
}
