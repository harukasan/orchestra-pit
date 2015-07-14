package exec_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/state/exec"
)

func TestExec(t *testing.T) {
	s := &exec.Exec{
		Script:     "echo apply",
		TestScript: "echo test",
	}

	if err := s.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("Test: %v", err)
	}
}
