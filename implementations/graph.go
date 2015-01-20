package implementations

import (
	"errors"
	"fmt"
	"github.com/twmb/algoimpl/go/graph"
	"reflect"
)

type graphType struct {
	graph  *graph.Graph
	phases map[string]graph.Node
}

var Graph = &graphType{
	graph:  graph.New(graph.Directed),
	phases: map[string]graph.Node{},
}

func (g *graphType) AddPhase(id string) error {
	node := g.graph.MakeNode()
	*node.Value = id
	g.phases[id] = node
	return nil
}

func (g *graphType) AddDependency(a, b string) error {
	va, found := g.phases[a]
	if !found {
		return fmt.Errorf("Could not find phase with id = %s", a)
	}
	vb, found := g.phases[b]
	if !found {
		return fmt.Errorf("Could not find phase with id = %s", b)
	}
	g.graph.MakeEdge(va, vb)
	return nil
}

func (g *graphType) Linearize() ([]string, error) {
	nodes := g.graph.TopologicalSort()
	ids := []string{}
	for _, n := range nodes {
		id, ok := (*n.Value).(string)
		if !ok {
			msg := fmt.Sprintf("Could not convert value: %v to string!", n.Value)
			if n.Value != nil {
				typ := reflect.TypeOf(*n.Value)
				msg += fmt.Sprintf(" Had type: %s", typ.String())
			}
			return nil, errors.New(msg)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (g *graphType) Reset() {
	g.graph = graph.New(graph.Directed)
	g.phases = map[string]graph.Node{}
}

func (g *graphType) String() string {
	return "Graph implementation"
}
