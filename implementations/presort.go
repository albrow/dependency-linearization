package implementations

import (
	"container/list"
	"fmt"
)

type presortType struct {
	phases   *list.List
	hasCycle bool
}

var Presort = &presortType{
	phases: list.New(),
}

type presortPhase struct {
	id   string
	deps []*presortPhase
}

func (p *presortType) AddPhase(id string) error {
	// Phases without any dependencies go in front
	p.phases.PushBack(&presortPhase{id: id})
	return nil
}

func (t *presortType) AddDependency(depId, pId string) error {
	if t.hasCycle {
		return nil
	}

	var p, dep *presortPhase
	var pEl, depEl *list.Element
	inBetweens := []*presortPhase{}
	for e := t.phases.Front(); e != nil; e = e.Next() {
		currentPhase, ok := e.Value.(*presortPhase)
		if !ok {
			return fmt.Errorf("Could not convert %v of type %T to *presortPhase!", e.Value, e.Value)
		}
		switch currentPhase.id {
		case pId:
			// We found p
			if dep != nil {
				// If we already found dep, that means dep came before p
				// in the list and the dependency is already satisfied. We
				// don't need to change the order.
				p = currentPhase
				p.deps = append(p.deps, dep)
				return nil
			}
			// If we're here it means we haven't found dep yet. The dependency
			// may be satisfiable but we'll have to move things around. We set
			// p and pEl here, so when we find dep we can figure out what to move
			// and where to move it. If we never find dep, we'll reach the end of
			// the function and return an error.
			p = currentPhase
			pEl = e
		case depId:
			// We found dep
			dep = currentPhase
			depEl = e
			if p != nil {
				// We found p but it was not before dep in the list. That means we'll
				// need to move some things around.
				p.deps = append(p.deps, dep)

				// If p depends on dep and dep depends on p, we have a pretty clear cycle
				for _, depdep := range dep.deps {
					if depdep.id == pId {
						t.hasCycle = true
						return nil
					}
				}

				// First, we'll attmpt to move p immediately after dep
				// We need to check any elements between p and dep
				// to see if they depend on p
				if anyDependsOn(inBetweens, p) {
				} else {
					// If we reached here, we found a placement that works!
					// We can move p directly after dep
					t.phases.MoveAfter(pEl, depEl)
					return nil
				}

				// Next, we'll attempt to move dep immediately before p
				// We need to make sure dep doesn't depend on any of the phases
				// between p and dep
				if dependsOnAny(dep, inBetweens) {
					// If we've reached here, we cannot move p direcly after dep
					// or dep directly before p, so we must have a cycle.
					t.hasCycle = true
					return nil
				} else {
					// If we reached here, we found a placement that works!
					// We can move dep directly before p
					t.phases.MoveBefore(depEl, pEl)
					return nil
				}
			}
		default:
			if p != nil {
				inBetweens = append(inBetweens, currentPhase)
			}
		}
	}

	if p == nil {
		return fmt.Errorf("Could not find phase with id = %s", pId)
	}
	if dep == nil {
		return fmt.Errorf("Could not find phase with id = %s", depId)
	}
	return nil
}

func (c *presortType) Linearize() ([]string, error) {
	// NOTE: if we can return a linked list here instead of a slice
	// of strings it would be even faster
	if c.hasCycle {
		return nil, fmt.Errorf("There was a cycle")
	}
	results := []string{}
	for e := c.phases.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(*presortPhase)
		if !ok {
			return nil, fmt.Errorf("Could not convert %v of type %T to *presortPhase!", e.Value, e.Value)
		}
		results = append(results, p.id)
	}
	return results, nil
}

// anyDependsOn returns true iff any phase in phases depends on p
func anyDependsOn(phases []*presortPhase, p *presortPhase) bool {
	for _, phase := range phases {
		for _, dep := range phase.deps {
			if dep.id == p.id {
				return true
			}
		}
	}
	return false
}

// dependsOnAny returns true iff phase depends on any phases
func dependsOnAny(p *presortPhase, phases []*presortPhase) bool {
	for _, phase := range phases {
		for _, dep := range p.deps {
			if dep.id == phase.id {
				return true
			}
		}
	}
	return false
}

func (t *presortType) phaseIds() []string {
	ids := []string{}
	for e := t.phases.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(*presortPhase)
		if !ok {
			msg := fmt.Sprintf("Could not convert %v of type %T to *presortPhase!", e.Value, e.Value)
			panic(msg)
		}
		ids = append(ids, p.id)
	}
	return ids
}

func phaseIds(phases []*presortPhase) []string {
	ids := []string{}
	for _, phase := range phases {
		ids = append(ids, phase.id)
	}
	return ids
}

func (p *presortPhase) depIds() []string {
	ids := []string{}
	for _, dep := range p.deps {
		ids = append(ids, dep.id)
	}
	return ids
}

func (c *presortType) Reset() {
	c.phases.Init()
}

func (c *presortType) String() string {
	return "Presort implementation"
}
