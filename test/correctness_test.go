package test

import (
	"github.com/albrow/dependency-linearization/common"
	"github.com/albrow/dependency-linearization/implementations"
	"testing"
)

var testCases = []struct {
	deps     []dep
	expected []string
}{
	{
		deps:     []dep{{"a", "b"}, {"b", "c"}},
		expected: []string{"a", "b", "c"},
	},
}

func TestGoraphGsKahn(t *testing.T) {
	testLinearizer(t, implementations.GoraphGsKahn)
}

func testLinearizer(t *testing.T, l common.Linearizer) {
	for _, tc := range testCases {
		runTestCase(t, l, tc.deps, tc.expected)
	}
}
