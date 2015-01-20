package common

type Linearizer interface {
	AddPhase(Phase) error
	Linearize() ([]Phase, error)
	// AddDependency adds b as a dependency to a.
	// It reads naturally as "a depends on b"
	AddDependency(a, b Phase) error
}
