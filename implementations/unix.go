package implementations

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type unixType struct {
	phases map[string]struct{}
	deps   []dep
}

type dep struct {
	depender  string
	dependsOn string
}

var Unix = &unixType{
	phases: map[string]struct{}{},
}

func (u *unixType) AddPhase(id string) error {
	u.phases[id] = struct{}{}
	return nil
}

func (u *unixType) AddDependency(a, b string) error {
	u.deps = append(u.deps, dep{a, b})
	return nil
}

func (u *unixType) Linearize() ([]string, error) {
	// Set up the tsort command and get the stdin pipe
	cmd := exec.Command("tsort")
	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// All the phases that aren't involved in any dependencies
	// need to be added later.
	leftOverPhases := map[string]struct{}{}
	// Copy all phases into leftOverPhases
	for p := range u.phases {
		leftOverPhases[p] = struct{}{}
	}

	// Write all the dependency statements to stdin
	for _, d := range u.deps {
		stmt := fmt.Sprintf("%s %s ", d.depender, d.dependsOn)
		stdin.Write([]byte(stmt))
		// As we go, delete the phases that are in some dependency
		// from the leftOverPhases. That way the only once left will
		// be those that are not involved in any dependencies
		delete(leftOverPhases, d.depender)
		delete(leftOverPhases, d.dependsOn)
	}
	// Write all the left over phases to stdin
	for p := range leftOverPhases {
		stmt := fmt.Sprintf("%s %s ", p, p)
		stdin.Write([]byte(stmt))
	}

	// Close pipe and wait for command to finish
	if err := stdin.Close(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// tsort reported an error
			return nil, fmt.Errorf("Error in tsort command: %s", out.String())
		} else {
			// There was some other problem
			return nil, err
		}
	}
	// If tsort worked, the output is one id per line
	// and the last line is an empty space. We can parse
	// that into a slice of ids pretty easily.
	ids := strings.Split(strings.TrimSpace(out.String()), "\n")
	return ids, nil
}

func (u *unixType) Reset() {
	u.deps = []dep{}
	u.phases = map[string]struct{}{}
}

func (u *unixType) String() string {
	return "Unix (builtin tsort) implementation"
}
