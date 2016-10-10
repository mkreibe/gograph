package graph

import (
	"reflect"
	"testing"
)

func printError(t *testing.T, testName string, expected interface{}, actual interface{}) {
	t.Error(
		"Test: ", testName,
		"Got["+reflect.TypeOf(actual).String()+"]: ", actual,
		"Expected["+reflect.TypeOf(expected).String()+"]: ", expected,
	)
}

func AssertWithCheckFunc(t *testing.T, testName string, expected interface{}, actual interface{}, checker func(e interface{}, a interface{}) bool) bool {
	if !checker(expected, actual) {
		printError(t, testName, actual, expected)
		return false
	}

	t.Log("Test Passed: ", testName)
	return true
}

func AssertT(t *testing.T, testName string, expected interface{}, actual interface{}) bool {
	return AssertWithCheckFunc(t, testName, expected, actual, func(e interface{}, a interface{}) bool {
		return e == a
	})
}

func AssertStrInterMap(t *testing.T, testName string, expected *map[string]interface{}, actual *map[string]interface{}) (result bool) {
	if result = AssertT(t, testName+"[count]", len(*expected), len(*actual)); result {
		for key, value := range *expected {
			thisTest := testName + "[item: " + key + "]"
			if v, ok := (*actual)[key]; ok {
				result = AssertT(t, thisTest, value, v)
			} else {
				printError(t, thisTest, key, "<nothing>")
				result = false
			}

			if !result {
				break
			}
		}
	}

	return
}

func AssertAttributes(t *testing.T, testName string, expected AttributeCollection, actual AttributeCollection) (result bool) {
	if result = AssertT(t, testName+"[count]", expected.Count(), actual.Count()); result {
		for key, value := range expected {
			thisTest := testName + "[item: " + key + "]"
			if v, ok := actual.Get(key); ok {
				result = AssertT(t, thisTest, value, v)
			} else {
				printError(t, thisTest, key, "<nothing>")
				result = false
			}

			if !result {
				break
			}
		}
	}

	return
}

func AssertEdges(t *testing.T, testName string, expected []*Edge, actual map[string]*Edge) (result bool) {
	if result = AssertT(t, testName+"[count]", len(expected), len(actual)); result {
		for _, edge := range expected {
			if result = AssertAttributes(t, testName+"[attrs]", edge.Attributes, actual[edge.ID].Attributes); !result {
				break
			}
		}
	}
	return
}

func AssertNodes(t *testing.T, testName string, expected []*Node, actual map[string]*Node) (result bool) {
	if result = AssertT(t, testName+"[count]", len(expected), len(actual)); result {
		for _, node := range expected {
			if result = AssertAttributes(t, testName+"[attrs]", node.Attributes, actual[node.ID].Attributes); !result {
				break
			}
		}
	}

	return
}

func AssertAreConnected(t *testing.T, graph *Graph, source string, target string, expected bool) (result bool) {

	symbol := "-"
	if graph.Type == GraphDirected {
		symbol = "->"
	}

	if actual, err := graph.HasConnection(source, target); err == nil {
		result = (expected == actual)

		if !result {
			printError(t, "Edge Check["+source+symbol+target+"]", expected, actual)
		}
	} else {
		t.Error(err)
	}
	return
}
