package implementations

import (
	"errors"
	"github.com/albrow/dependency-linearization/common"
	"github.com/gyuho/goraph/algorithm/tskahn"
	"github.com/gyuho/goraph/graph/gs"
	"strings"
)

type goraphGsKahnType struct {
	graph *gs.Graph
}

var GoraphGsKahn = &goraphGsKahnType{
	graph: gs.NewGraph(),
}

func (g *goraphGsKahnType) AddPhase(p common.Phase) error {
	vertex := gs.NewVertex(p.Id())
	g.graph.AddVertex(vertex)
	return nil
}

func (g *goraphGsKahnType) AddDependency(a, b common.Phase) error {
	va := g.graph.FindVertexByID(a.Id())
	vb := g.graph.FindVertexByID(b.Id())
	g.graph.Connect(va, vb, 0)
	return nil
}

func (g *goraphGsKahnType) Linearize() ([]common.Phase, error) {
	// TODO: actually sort this with some algorithm
	sorted, ok := tskahn.TSKahn(g.graph)
	if !ok {
		return nil, errors.New("Could not linearize dependencies. Was there a cycle?")
	}
	ids := strings.Split(sorted, " → ")
	phases := []common.Phase{}
	for _, id := range ids {
		phases = append(phases, common.NewPhase(id))
	}
	return phases, nil
}

func (g *goraphGsKahnType) String() string {
	return "Goraph with gs type and Khan algorithm"
}
