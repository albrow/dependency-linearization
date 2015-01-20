package common

type Linearizer interface {
	AddPhase(string) error
	Linearize() ([]string, error)
	// AddDependency adds b as a dependency to a.
	// It reads naturally as "a depends on b"
	AddDependency(a, b string) error
	// Reset clears all previous phases
	Reset()
}
