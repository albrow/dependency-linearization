package test

import (
	"github.com/albrow/dependency-linearization/common"
	"github.com/albrow/dependency-linearization/implementations"
	"testing"
)

var (
	linear1Deps  = makeLinearDeps(1)
	linear3Deps  = makeLinearDeps(3)
	linear10Deps = makeLinearDeps(10)
	tree1Deps    = makeTreeDeps(1)
	tree3Deps    = makeTreeDeps(3)
	tree10Deps   = makeTreeDeps(10)
)

func BenchmarkLinear1Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, linear1Deps)
}

func BenchmarkLinear1Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, linear1Deps)
}

func BenchmarkLinear1Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, linear1Deps)
}

func BenchmarkLinear1Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, linear1Deps)
}

func BenchmarkLinear3Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, linear3Deps)
}

func BenchmarkLinear3Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, linear3Deps)
}

func BenchmarkLinear3Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, linear3Deps)
}

func BenchmarkLinear3Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, linear3Deps)
}

func BenchmarkLinear10Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, linear10Deps)
}

func BenchmarkLinear10Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, linear10Deps)
}

func BenchmarkLinear10Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, linear10Deps)
}

func BenchmarkLinear10Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, linear10Deps)
}

func BenchmarkTree1Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, tree1Deps)
}

func BenchmarkTree1Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, tree1Deps)
}

func BenchmarkTree1Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, tree1Deps)
}

func BenchmarkTree1Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, tree1Deps)
}

func BenchmarkTree3Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, tree3Deps)
}

func BenchmarkTree3Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, tree3Deps)
}

func BenchmarkTree3Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, tree3Deps)
}

func BenchmarkTree3Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, tree3Deps)
}

func BenchmarkTree10Goraph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Goraph, tree10Deps)
}

func BenchmarkTree10Unix(b *testing.B) {
	benchmarkLinearizer(b, implementations.Unix, tree10Deps)
}

func BenchmarkTree10Graph(b *testing.B) {
	benchmarkLinearizer(b, implementations.Graph, tree10Deps)
}

func BenchmarkTree10Maps(b *testing.B) {
	benchmarkLinearizer(b, implementations.Maps, tree10Deps)
}

// benchmarkLinearizer runs the given deps list through
// the linearizer and benchmarks the time it takes to 1) add each phase,
// 2) add each dependency, and 3) linearize. It attempts to do so with
// as little overhead as possible by using a preparer and stopping the timer
// to do error checking code and bookkeeping.
func benchmarkLinearizer(b *testing.B, l common.Linearizer, deps []dep) {
	p := prepareCase(l, deps)
	funcs := p.getFuncs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, f := range funcs {
			err := f()
			b.StopTimer()
			if err != nil {
				panic(err)
			}
			b.StartTimer()
		}
		_, err := l.Linearize()
		b.StopTimer()
		if err != nil {
			panic(err)
		}
		l.Reset()
		b.StartTimer()
	}
}
