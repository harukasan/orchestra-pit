// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build darwin

package packagemanager_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/state/packagemanager"
)

func TestInstalled(t *testing.T) {
	s := &packagemanager.Installed{
		Name:    "orchestra-pit/fake/stub",
		Options: []string{"--with-something-awesome"},
		Update:  true,
	}

	if err := s.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}
	if err := s.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("Test: %v", err)
	}
}

func TestRemoved(t *testing.T) {
	is := &packagemanager.Installed{
		Name: "orchestra-pit/fake/stub",
	}
	if err := is.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}

	s := &packagemanager.Removed{
		Name: "orchestra-pit/fake/stub",
	}
	if err := s.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("Test: %v", err)
	}
}
