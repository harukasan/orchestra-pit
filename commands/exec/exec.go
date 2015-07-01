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

// Command returns the initialized Cmd struct to execute the named program with
// the given arguments. It implements exec.Cmd.
func Command(name string, arg ...string) *Cmd {
	c := exec.Command(name, arg...)
	return &Cmd{c}
}

// Exec is a dummy state which executes the given script when the state is
// applied.
//
// Script specifies the command to execute on that the state is applied. It
// executes on applying the state. If the command exits with non-zero status,
// Apply returns an error.
//
// TestScript specifies the command to execute on testing the state. If the
// command exits with non-zero status, Test returns an error.
//
type Exec struct {
	Script     string
	TestScript string
	Env        []string
}

// Apply tries to execute the given script. If the script exits with non-zero
// status, Apply returns an error.
func (s *Exec) Apply() error {
	args := strings.Split(s.Script, " ")
	c := Command(args[0], args[1:]...)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, s.Env...)
	return c.Run()
}

// Test tries to execute the given TestScript, if the TestScript parameter is
// not empty. If the script exits with non-zero status, Test returns an error.
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
