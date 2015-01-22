package implementations

import (
	"container/list"
	"fmt"
	"strings"
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
	deps []string
}

func (p *presortType) AddPhase(id string) error {
	// Phases without any dependencies go in front
	p.phases.PushBack(&presortPhase{id: id})
	return nil
}

func (c *presortType) AddDependency(a, b string) error {
	if c.hasCycle {
		return nil
	}

	// First find the elements that correspond to a and b
	var aEl *list.Element
	var bEl *list.Element
	var bp *presortPhase
	seen := ""
	for e := c.phases.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(*presortPhase)
		if !ok {
			return fmt.Errorf("Could not convert %v of type %T to *presortPhase!", e.Value, e.Value)
		}
		switch p.id {
		case a:
			aEl = e
			if bEl != nil {
				// Check if we have already seen all the dependencies for
				// b. If we have not, then there is a cycle.
				for _, dep := range bp.deps {
					if !strings.Contains(seen, dep+" ") {
						c.hasCycle = true
						return nil
					}
				}
				// Move b to after a
				c.phases.MoveAfter(bEl, aEl)
				return nil
			}
		case b:
			if aEl != nil {
				// If this is the case, b was already after a.
				// We don't need to move anything
				return nil
			}
			p.deps = append(p.deps, a)
			bp = p
			bEl = e
		}
		seen += p.id + " "
	}
	if aEl == nil {
		return fmt.Errorf("Could not find phase with id = %s", a)
	}
	if bEl == nil {
		return fmt.Errorf("Could not find phase with id = %s", b)
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

func (c *presortType) Reset() {
	c.phases.Init()
}

func (c *presortType) String() string {
	return "Presort implementation"
}
