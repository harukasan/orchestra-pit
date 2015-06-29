package exec

import "os/exec"

// Cmd represents an external command.
type Cmd struct {
	*exec.Cmd
}

// Command returns the Cmd struct to execute the named program with
// the given arguments. It implements exec.Cmd.
func Command(name string, arg ...string) *Cmd {
	c := exec.Command(name, arg...)
	return &Cmd{c}
}
