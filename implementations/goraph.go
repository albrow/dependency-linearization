package implementations

import (
	"errors"
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

func (g *goraphGsKahnType) AddPhase(id string) error {
	vertex := gs.NewVertex(id)
	g.graph.AddVertex(vertex)
	return nil
}

func (g *goraphGsKahnType) AddDependency(a, b string) error {
	va := g.graph.FindVertexByID(a)
	vb := g.graph.FindVertexByID(b)
	g.graph.Connect(va, vb, 0)
	return nil
}

func (g *goraphGsKahnType) Linearize() ([]string, error) {
	// TODO: actually sort this with some algorithm
	sorted, ok := tskahn.TSKahn(g.graph)
	if !ok {
		return nil, errors.New("Could not linearize dependencies. Was there a cycle?")
	}
	ids := strings.Split(sorted, " → ")
	return ids, nil
}

func (g *goraphGsKahnType) Reset() {
	g.graph = gs.NewGraph()
}

func (g *goraphGsKahnType) String() string {
	return "Goraph with gs type and Khan algorithm"
}
