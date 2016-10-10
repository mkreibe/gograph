package graph

import (
	"errors"
	"testing"
)

func TestBadGraphType(t *testing.T) {
	graph := Graph{Type: "foo"}
	if _, err := graph.HasConnection("first", "second"); err != nil {
		errMsg := "Unknown graph type: foo"
		if err.Error() != errMsg {
			printError(t, "Bad graph type", errMsg, err.Error())
		}
	} else {
		t.Error("Expected an error, invalid type.")
	}
}

func TestBadNodeId(t *testing.T) {
	if _, err := LoadFileGraph("./data/e_nonodeid.json"); err != nil {
		errMsg := "Node must have a valid id"
		if err.Error() != errMsg {
			printError(t, "Bad node id", errMsg, err.Error())
		}
	} else {
		t.Error("Expected an error, invalid node id.")
	}
}

func TestBadEdgeId(t *testing.T) {
	if _, err := LoadFileGraph("./data/e_noedgeid.json"); err != nil {
		errMsg := "Invalid edge id"
		if err.Error() != errMsg {
			printError(t, "Bad edge id", errMsg, err.Error())
		}
	} else {
		t.Error("Expected an error, invalid edge id.")
	}
}

func TestBadEdgeConnectionsNode(t *testing.T) {
	if _, err := LoadFileGraph("./data/e_noedgeconnections.json"); err != nil {
		errMsg := "Edges must have atleast two nodes to connect"
		if err.Error() != errMsg {
			printError(t, "Bad edge connection", errMsg, err.Error())
		}
	} else {
		t.Error("Expected an error, invalid edge connections (none).")
	}
}

func TestBadEdgeConnectionsSingular(t *testing.T) {
	if _, err := LoadFileGraph("./data/e_singleedgeconnection.json"); err != nil {
		errMsg := "Edges must have atleast two nodes to connect"
		if err.Error() != errMsg {
			printError(t, "Bad edge connection (single)", errMsg, err.Error())
		}
	} else {
		t.Error("Expected an error, invalid edge connections (none).")
	}
}

func TestGraphLoading(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {
		AssertT(t, "Type", GraphUndirected, graph.Type)

		attributes := NewAttributeCollection()
		attributes.Set("description", "simple 5 node, 7 edge graph. Taken from CRLS.")

		if !AssertAttributes(t, "Attrib", attributes, graph.Attributes) {
			return
		}

		nodes := []*Node{
			&Node{ID: "1", Attributes: BuildAttributes(map[string]interface{}{"color": "green"})},
			&Node{ID: "2"},
			&Node{ID: "3"},
			&Node{ID: "4"},
			&Node{ID: "5"},
		}

		if !AssertNodes(t, "Node", nodes, graph.Nodes) {
			return
		}

		edges := []*Edge{
			&Edge{ID: "a", Attributes: BuildAttributes(map[string]interface{}{"style": "dashed"})},
			&Edge{ID: "b"},
			&Edge{ID: "c"},
			&Edge{ID: "d"},
			&Edge{ID: "e"},
			&Edge{ID: "f"},
			&Edge{ID: "g"},
		}

		if !AssertEdges(t, "Edge", edges, graph.Edges) {
			return
		}

		// n - n connections are implied with undirected graphs.

		AssertAreConnected(t, graph, "1", "1", true)
		AssertAreConnected(t, graph, "1", "2", true)
		AssertAreConnected(t, graph, "2", "1", true)
		AssertAreConnected(t, graph, "1", "3", false)
		AssertAreConnected(t, graph, "3", "1", false)
		AssertAreConnected(t, graph, "1", "4", false)
		AssertAreConnected(t, graph, "4", "1", false)
		AssertAreConnected(t, graph, "1", "5", true)
		AssertAreConnected(t, graph, "5", "1", true)

		AssertAreConnected(t, graph, "2", "2", true)
		AssertAreConnected(t, graph, "2", "3", true)
		AssertAreConnected(t, graph, "3", "2", true)
		AssertAreConnected(t, graph, "2", "4", true)
		AssertAreConnected(t, graph, "4", "2", true)
		AssertAreConnected(t, graph, "2", "5", true)
		AssertAreConnected(t, graph, "5", "2", true)

		AssertAreConnected(t, graph, "3", "3", true)
		AssertAreConnected(t, graph, "3", "4", true)
		AssertAreConnected(t, graph, "4", "3", true)
		AssertAreConnected(t, graph, "5", "3", false)
		AssertAreConnected(t, graph, "3", "5", false)

		AssertAreConnected(t, graph, "4", "4", true)
		AssertAreConnected(t, graph, "4", "5", true)
		AssertAreConnected(t, graph, "5", "4", true)

		AssertAreConnected(t, graph, "5", "5", true)
	}
}

func TestIterInvalidAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		var invalidAttrValue *string

		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", invalidAttrValue)
		}

		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			return // never called.
		}); err != nil {
			errMsg := "Non sortable type. Requires string, int or float64"
			if err.Error() != errMsg {
				printError(t, "Bad invalid attr", errMsg, err.Error())
			}
		} else {
			t.Error("Expected an error. Iteration, no attr.")
		}
	}
}

func TestNoIterFunc(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {
		if err = graph.IterNodes("foo", nil); err != nil {
			errMsg := "No iteration function"
			if err.Error() != errMsg {
				printError(t, "Bad invalid attr", errMsg, err.Error())
			}
		} else {
			t.Error("Expected an error. Iteration, no attr.")
		}
	}
}

func TestIterStringAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		attrValue := "a"
		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", attrValue)
		}

		count := 0
		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			count++
			cont = true
			return
		}); err == nil {
			AssertT(t, "String attr", 5, count)
		} else {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestIterIntAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		attrValue := 1
		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", attrValue)
		}

		count := 0
		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			count++
			cont = true
			return
		}); err == nil {
			AssertT(t, "Int attr", 5, count)
		} else {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestIterFloatAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		attrValue := 0.124
		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", attrValue)
		}

		count := 0
		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			count++
			cont = true
			return
		}); err == nil {
			AssertT(t, "Int attr", 5, count)
		} else {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestIterBreakAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		attrValue := 1
		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", attrValue)
		}

		count := 0
		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			count++
			cont = count > 3
			return
		}); err == nil {
			AssertT(t, "break attr", 2, count)
		} else {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestIterErrAttr(t *testing.T) {
	if graph, err := LoadFileGraph("./data/ud_5n_7e.json"); err != nil {
		t.Error(err)
	} else {

		attrValue := 1
		for _, n := range graph.Nodes {
			n.Attributes.Set("foo", attrValue)
		}

		errMsg := "Invalid edge id"
		if err = graph.IterNodes("foo", func(node *Node) (cont bool, err error) {
			err = errors.New(errMsg)
			cont = true // make sure that if cont is true, that we still handle the break.
			return
		}); err != nil {
			if err.Error() != errMsg {
				printError(t, "Bad iteration break", errMsg, err.Error())
			}
		} else {
			t.Error("Expected an error. Iteration error.")
		}
	}
}
