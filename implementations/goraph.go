package implementations

import (
	"errors"
	"github.com/gyuho/goraph/algorithm/tsdag"
	"github.com/gyuho/goraph/graph/gs"
	"strings"
)

type goraphType struct {
	graph *gs.Graph
}

var Goraph = &goraphType{
	graph: gs.NewGraph(),
}

func (g *goraphType) AddPhase(id string) error {
	vertex := gs.NewVertex(id)
	g.graph.AddVertex(vertex)
	return nil
}

func (g *goraphType) AddDependency(a, b string) error {
	va := g.graph.FindVertexByID(a)
	vb := g.graph.FindVertexByID(b)
	g.graph.Connect(va, vb, 0)
	return nil
}

func (g *goraphType) Linearize() ([]string, error) {
	// TODO: actually sort this with some algorithm
	sorted, ok := tsdag.TSDAG(g.graph)
	if !ok {
		return nil, errors.New("Could not linearize dependencies. Was there a cycle?")
	}
	ids := strings.Split(sorted, " → ")
	return ids, nil
}

func (g *goraphType) Reset() {
	g.graph = gs.NewGraph()
}

func (g *goraphType) String() string {
	return "Goraph with gs type and Khan algorithm"
}
