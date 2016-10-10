package graph

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// BFSd is the breadth-first-search algorithm on the graph with the attribute set
// as 'd'.
func (graph *Graph) BFSd(root string) (err error) {
	return graph.BFS(root, "d")
}

// BFS is the breadth-first-search algorithm on the graph.
func (graph *Graph) BFS(root string, valueAttr string) (err error) {
	colorAttr := randStringRunes(48)
	predecessorAttr := randStringRunes(48)

	rootNode := graph.Nodes[root]

	for name, val := range graph.Nodes {
		if name != root {
			val.Attributes.Set(colorAttr, "white")
			val.Attributes.Set(valueAttr, -1) // Our version of infinity.
			val.Attributes.Set(predecessorAttr, nil)
		}
	}

	rootNode.Attributes.Set(colorAttr, "gray")
	rootNode.Attributes.Set(valueAttr, 0)
	rootNode.Attributes.Set(predecessorAttr, nil)

	var Q []*Node
	for Q = append(Q, rootNode); len(Q) != 0; {
		u := Q[0]
		Q = Q[1:]

		// for each adjacent
		if rawVal, ok := u.Attributes.Get(valueAttr); ok {
			value := rawVal.(int)
			value++
			for _, v := range u.Adj() {
				if color, ok := v.Attributes.Get(colorAttr); ok && color == "white" {
					v.Attributes.Set(colorAttr, "gray")
					v.Attributes.Set(valueAttr, value)
					v.Attributes.Set(predecessorAttr, u)
					Q = append(Q, v)
				}
			}
		}

		u.Attributes.Set(colorAttr, "black")
	}

	// Cleanup!
	for _, val := range graph.Nodes {
		val.Attributes.Remove(colorAttr)
		val.Attributes.Remove(predecessorAttr)
	}

	return
}
