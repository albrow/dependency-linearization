package implementations

import (
	"fmt"
)

type mapsType struct {
	// A map of phases to the phases they depend on
	phases map[string]map[string]struct{}
}

var Maps = &mapsType{
	phases: map[string]map[string]struct{}{},
}

func (c *mapsType) AddPhase(id string) error {
	c.phases[id] = map[string]struct{}{}
	return nil
}

func (c *mapsType) AddDependency(a, b string) error {
	if _, found := c.phases[b]; !found {
		return fmt.Errorf("Could not find phase with id = %s", b)
	}
	c.phases[b][a] = struct{}{}
	return nil
}

func (c *mapsType) Linearize() ([]string, error) {
	results := []string{}
	for len(c.phases) > 0 {
		deleted := []string{}
		for phase, deps := range c.phases {
			// Find the phases which have no dependencies left
			// Add them to results and remove them from deps map
			if len(deps) == 0 {
				results = append(results, phase)
				delete(c.phases, phase)
				deleted = append(deleted, phase)
			}
		}
		if len(deleted) == 0 {
			return nil, fmt.Errorf("Detected cycle!")
		} else {
			for _, phaseToDelete := range deleted {
				for _, deps := range c.phases {
					delete(deps, phaseToDelete)
				}
			}
		}
	}
	return results, nil
}

func mapKeys(m map[string]struct{}) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func (c *mapsType) Reset() {
	c.phases = map[string]map[string]struct{}{}
}

func (c *mapsType) String() string {
	return "Maps implementation"
}
