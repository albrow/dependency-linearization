package test

import (
	"github.com/albrow/dependency-linearization/common"
	"testing"
)

// dep defines a dependency by two
// phase ids. If you want to represent
// only one phase, leave dependsOn blank
type dep struct {
	depender  string
	dependsOn string
}

// runTestCase runs l against a specific test case, which is defined
// by deps. expected should be a slice of phase ids in the expected
// order.
func runTestCase(t *testing.T, l common.Linearizer, deps []dep, expected []string) {
	defer l.Reset()
	// Create the phases and set up dependencies as needed
	phases := newPhaseList(l)
	for _, dep := range deps {
		phases.nxAdd(dep.depender)
		if dep.dependsOn != "" {
			phases.nxAdd(dep.dependsOn)
			if err := l.AddDependency(dep.depender, dep.dependsOn); err != nil {
				panic(err)
			}
		}
	}
	got, err := l.Linearize()
	if err != nil {
		panic(err)
	}
	compareResults(t, l, got, expected)
}

func compareResults(t *testing.T, l common.Linearizer, got []string, expected []string) {
	if len(expected) != len(got) {
		t.Errorf("Results were not the correct length for %s. Expected %d phases but got %d.\n\tExpected: %v\n\tGot: %v",
			l, len(expected), len(got), expected, got)
		return
	}
	for i, e := range expected {
		g := got[i]
		if e != g {
			t.Errorf("Phase[%d] of results was incorrect for %s. Expected phase id = %s but got id = %s\n\tExpected: %v\n\tGot: %v",
				i, l, e, g, expected, got)
			return
		}
	}
}

// phaseList is used for managing a list of phases
// only really for testing purposes
type phaseList struct {
	phases     map[string]struct{}
	linearizer common.Linearizer
}

func newPhaseList(l common.Linearizer) *phaseList {
	return &phaseList{
		phases:     map[string]struct{}{},
		linearizer: l,
	}
}

// nxAdd only adds the id if it doesn't already exist. It will
// not add duplicate ids to the phaseList
func (pl *phaseList) nxAdd(id string) {
	if _, found := pl.phases[id]; !found {
		// only add the phase to the phaseList and
		// the underlying graph if it hasn't already
		// been added
		pl.phases[id] = struct{}{}
		if err := pl.linearizer.AddPhase(id); err != nil {
			panic(err)
		}
	}
}
