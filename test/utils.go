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
	// Create the phases and set up dependencies as needed
	phases := newPhaseList(l)
	for _, dep := range deps {
		p := phases.getOrCreate(dep.depender)
		if dep.dependsOn != "" {
			d := phases.getOrCreate(dep.dependsOn)
			if err := l.AddDependency(p, d); err != nil {
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

func compareResults(t *testing.T, l common.Linearizer, got []common.Phase, expected []string) {
	if len(expected) != len(got) {
		t.Errorf("Results were not the correct length for %s. Expected %d phases but got %d.\n\tExpected: %v\n\tGot: %v",
			l, len(expected), len(got), expected, getIdsForPhases(got))
		return
	}
	for i, e := range expected {
		g := got[i]
		if e != g.Id() {
			t.Errorf("Phase[%d] of results was incorrect for %s. Expected phase id = %s but got id = %s\n\tExpected: %v\n\tGot: %v",
				i, l, e, g, expected, getIdsForPhases(got))
			return
		}
	}
}

// phaseList is used for managing a list of phases
// only really for testing purposes
type phaseList struct {
	phases     map[string]common.Phase
	linearizer common.Linearizer
}

func newPhaseList(l common.Linearizer) *phaseList {
	return &phaseList{
		phases:     map[string]common.Phase{},
		linearizer: l,
	}
}

func (pl *phaseList) getOrCreate(id string) common.Phase {
	// if a phase with that id already exists,
	// return it
	if p, found := pl.phases[id]; found {
		return p
	}
	// otherwise create and return a new phase
	p := common.NewPhase(id)
	pl.phases[id] = p
	if err := pl.linearizer.AddPhase(p); err != nil {
		panic(err)
	}
	return p
}

func getIdsForPhases(phases []common.Phase) []string {
	ids := []string{}
	for _, p := range phases {
		ids = append(ids, p.Id())
	}
	return ids
}
