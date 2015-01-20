package test

import (
	"github.com/albrow/dependency-linearization/common"
	"strconv"
	"testing"
)

type testCase struct {
	deps     []dep
	expected []string
}

// dep defines a dependency by two
// phase ids. If you want to represent
// only one phase, leave dependsOn blank
type dep struct {
	depender  string
	dependsOn string
}

// runTestCase runs l against a specific test case, which is defined
// by tc.deps. expected should be a slice of phase ids in the expected
// order.
func runTestCase(t *testing.T, l common.Linearizer, tc testCase) {
	defer l.Reset()
	// Create the phases and set up dependencies as needed
	if err := prepareCase(l, tc.deps).execute(); err != nil {
		panic(err)
	}
	got, err := l.Linearize()
	if err != nil {
		panic(err)
	}
	compareResults(t, l, got, tc.expected)
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

// preparer does some bookkeeping for a list of phases
// and dependencies associated with some test case.
// After everything is prepared, you can call execute to
// have the preparer add the needed phases and dependencies
// to the given linearizer. It works by doing all the bookkeeping
// and preparation ahead of time, creating a series of functions
// to be run at a later time with minimal overhead.
// That way you can benchmark only the function calls
// themselves and not the prep work which went into
// organizing them.
type preparer struct {
	addedPhases map[string]struct{}
	l           common.Linearizer
	phaseFuncs  []func() error
	depFuncs    []func() error
}

// prepareCase creates a new preparer with the given linearizer
// and deps. It prepares all the needed functions. All you
// have lef to do is call execute.
func prepareCase(l common.Linearizer, deps []dep) *preparer {
	p := &preparer{
		addedPhases: map[string]struct{}{},
		l:           l,
	}
	for _, dep := range deps {
		p.nxAddPhase(dep.depender)
		if dep.dependsOn != "" {
			p.nxAddPhase(dep.dependsOn)
			p.addDep(dep)
		}
	}
	return p
}

// nxAddPhase adds a function to the preparer which will add a phase with
// the given id to the preparer's linearizer iff such a function has not
// already been added. When execute is called, all such functions will be
// executed and each phase will be added exactly once.
func (p *preparer) nxAddPhase(id string) {
	if _, found := p.addedPhases[id]; !found {
		p.addedPhases[id] = struct{}{}
		p.phaseFuncs = append(p.phaseFuncs, func() error {
			return p.l.AddPhase(id)
		})
	}
}

// addDep adds a function to the preparer which will add the given dependency
// to the preparer's linearizer. When execute is called, all such functions
// will be executed.
func (p *preparer) addDep(dep dep) {
	p.depFuncs = append(p.depFuncs, func() error {
		return p.l.AddDependency(dep.depender, dep.dependsOn)
	})
}

// execute calls all the prepared functions in order
// It returns any errors that may have resulted. For finer
// control (e.g. to stop and start a benchmark timer around
// error checking code), you can use getFuncs to get the
// combined list of functions direcly.
func (p *preparer) execute() error {
	for _, addPhase := range p.phaseFuncs {
		if err := addPhase(); err != nil {
			return err
		}
	}
	for _, addDep := range p.depFuncs {
		if err := addDep(); err != nil {
			return err
		}
	}
	return nil
}

// getFuncs returns a combined list of all the functions
// the preparer has prepared, in the correct order. You
// can use it to gain finer control over timing in benchmarks.
func (p *preparer) getFuncs() []func() error {
	allFuncs := []func() error{}
	for _, f := range p.phaseFuncs {
		allFuncs = append(allFuncs, f)
	}
	for _, f := range p.depFuncs {
		allFuncs = append(allFuncs, f)
	}
	return allFuncs
}

// makeTreeDeps returns a slice of deps arranged in a tree pattern.
// i.e. some root phase depends on numBranches phases.
// Like this:
//
//       a
//     / | \
//    b  c  d
//
func makeTreeDeps(numBranches int) []dep {
	if numBranches == 0 {
		return []dep{{"0", ""}}
	}
	deps := []dep{}
	// start with 1 and iterate to numBranches
	for i := 1; i <= numBranches; i++ {
		deps = append(deps, dep{
			depender:  "0",
			dependsOn: strconv.Itoa(i),
		})
	}
	return deps
}

// makeLinearDeps returns a slice of numPhases deps arranged in a linear pattern.
// Like this:
//
//   a -> b -> c -> d
//
func makeLinearDeps(numPhases int) []dep {
	if numPhases == 0 {
		return []dep{{"0", ""}}
	}
	deps := []dep{}
	// start with 1 and iterate to numPhases
	for i := 1; i < numPhases; i++ {
		deps = append(deps, dep{
			depender:  strconv.Itoa(i - 1),
			dependsOn: strconv.Itoa(i),
		})
	}
	return deps
}
