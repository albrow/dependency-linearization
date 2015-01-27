package test

import (
	"github.com/albrow/dependency-linearization/common"
	"github.com/albrow/dependency-linearization/implementations"
	"strings"
	"testing"
)

var testCases = []testCase{
	{
		deps:     []dep{{"a", ""}},
		expected: []string{"a"},
	},
	{
		deps:     []dep{{"a", "b"}, {"b", "c"}},
		expected: []string{"a", "b", "c"},
	},
	{
		deps:     []dep{{"a", "b"}, {"b", "c"}, {"c", "d"}},
		expected: []string{"a", "b", "c", "d"},
	},
	{
		deps:     []dep{{"a", "b"}, {"b", "c"}, {"c", "d"}, {"a", "d"}},
		expected: []string{"a", "b", "c", "d"},
	},
}

func TestGoraph(t *testing.T) {
	testLinearizer(t, implementations.Goraph)
}

func TestGoraphCycle(t *testing.T) {
	testCycle(t, implementations.Goraph)
}

func TestUnix(t *testing.T) {
	testLinearizer(t, implementations.Unix)
}

func TestUnixCycle(t *testing.T) {
	testCycle(t, implementations.Unix)
}

func TestGraph(t *testing.T) {
	testLinearizer(t, implementations.Graph)
}

func TestGraphCycle(t *testing.T) {
	testCycle(t, implementations.Graph)
}

func TestMaps(t *testing.T) {
	testLinearizer(t, implementations.Maps)
}

func TestMapsCycle(t *testing.T) {
	testCycle(t, implementations.Maps)
}

func TestLists(t *testing.T) {
	testLinearizer(t, implementations.Lists)
}

func TestListsCycle(t *testing.T) {
	testCycle(t, implementations.Lists)
}

func TestPresort(t *testing.T) {
	testLinearizer(t, implementations.Presort)
}

func TestPresortCycle(t *testing.T) {
	testCycle(t, implementations.Presort)
}

func testLinearizer(t *testing.T, l common.Linearizer) {
	for _, tc := range testCases {
		runTestCase(t, l, tc)
	}
}

func testCycle(t *testing.T, l common.Linearizer) {
	// There is a cycle a -> b -> c -> a
	deps := []dep{{"a", "b"}, {"b", "c"}, {"c", "a"}, {"b", "d"}, {"d", "e"}}
	if err := prepareCase(l, deps).execute(); err != nil {
		panic(err)
	}
	if _, err := l.Linearize(); err == nil {
		t.Error("Expected error for cyclical graph but got none")
	} else {
		if !strings.Contains(err.Error(), "cycle") {
			t.Errorf("Expected error to say something about a cycle but it did not. Got: %s", err.Error())
		}
	}
	l.Reset()
}
