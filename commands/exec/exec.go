// Copyright 2015 MICHII Shunsuke. All rights reserved.

/*
Package exec provides the execution of commands.
*/
package exec

import (
	"os"
	"os/exec"
	"strings"
)

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

// Exec is a dummy state which executes the given script when the state is
// applied.
//
// Script specifies the command to execute on that the state is applied.
//
// TestScript specifies the command to execute on testing the state. If the
// command returns an errro when the command exits with non-zero status.
//
type Exec struct {
	Script     string
	TestScript string
	Env        []string
}

// Apply tries to execute the given script.
func (s *Exec) Apply() error {
	args := strings.Split(s.Script, " ")
	c := Command(args[0], args[1:]...)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, s.Env...)
	return c.Run()
}

// Test tries to execute the given TestScript, if the TestScript parameter is
// not empty.
func (s *Exec) Test() error {
	if s.TestScript == "" {
		return nil
	}

	args := strings.Split(s.TestScript, " ")
	c := Command(args[0], args[1:]...)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, s.Env...)
	return c.Run()
}
