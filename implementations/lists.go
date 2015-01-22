package implementations

import (
	"container/list"
	"fmt"
)

type listsType struct {
	// A linked list of linked lists representing dependencies
	phases *list.List
}

var Lists = &listsType{
	phases: list.New(),
}

type phase struct {
	deps *list.List
	id   string
}

func (l *listsType) AddPhase(id string) error {
	l.phases.PushFront(phase{
		deps: list.New(),
		id:   id,
	})
	return nil
}

func (c *listsType) AddDependency(a, b string) error {
	for e := c.phases.Front(); e != nil; e = e.Next() {
		p, ok := e.Value.(phase)
		if !ok {
			return fmt.Errorf("Could not convert %v of type %T to phase!", e.Value, e.Value)
		}
		if p.id == b {
			p.deps.PushFront(a)
			return nil
		}
	}
	return fmt.Errorf("Could not find phase with id = %s", b)
}

func (c *listsType) Linearize() ([]string, error) {
	results := []string{}
	for c.phases.Len() > 0 {
		deleted := []string{}
		for e := c.phases.Front(); e != nil; e = e.Next() {
			p, ok := e.Value.(phase)
			if !ok {
				return nil, fmt.Errorf("Could not convert %v of type %T to phase!", e.Value, e.Value)
			}
			if p.deps.Len() == 0 {
				results = append(results, p.id)
				c.phases.Remove(e)
				deleted = append(deleted, p.id)
			}
		}
		if len(deleted) == 0 {
			return nil, fmt.Errorf("Detected cycle!")
		} else {
			for _, phaseToDelete := range deleted {
				for e := c.phases.Front(); e != nil; e = e.Next() {
					p, ok := e.Value.(phase)
					if !ok {
						return nil, fmt.Errorf("Could not convert %v of type %T to phase!", e.Value, e.Value)
					}
					for dep := p.deps.Front(); dep != nil; dep = dep.Next() {
						depId, ok := dep.Value.(string)
						if !ok {
							return nil, fmt.Errorf("Could not convert %v of type %T to string!", dep.Value, dep.Value)
						}
						if depId == phaseToDelete {
							p.deps.Remove(dep)
						}
					}
				}
			}
		}
	}
	return results, nil
}

func (c *listsType) Reset() {
	c.phases.Init()
}

func (c *listsType) String() string {
	return "Lists implementation"
}
