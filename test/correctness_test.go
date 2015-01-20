package test

import (
	"github.com/albrow/dependency-linearization/common"
	"github.com/albrow/dependency-linearization/implementations"
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
}

func TestGoraph(t *testing.T) {
	testLinearizer(t, implementations.Goraph)
}

func TestUnix(t *testing.T) {
	testLinearizer(t, implementations.Unix)
}

func TestGraph(t *testing.T) {
	testLinearizer(t, implementations.Graph)
}

func testLinearizer(t *testing.T, l common.Linearizer) {
	for _, tc := range testCases {
		runTestCase(t, l, tc)
	}
}
