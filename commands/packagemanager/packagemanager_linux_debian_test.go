// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package packagemanager_test

import (
	"os"
	"testing"

	"github.com/harukasan/orchestra-pit/commands/packagemanager"
	"github.com/harukasan/orchestra-pit/commands/platform"
)

func TestMain(m *testing.M) {
	p, err := platform.Identify()
	if err != nil {
		return
	}
	if p.Get("family") != string(platform.FamilyDebian) {
		return
	}

	os.Exit(m.Run())
}

func TestInstalled(t *testing.T) {
	s := &packagemanager.Installed{
		Name: "debian-faq",
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
		Name: "debian-faq",
	}
	if err := is.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}

	s := &packagemanager.Removed{
		Name: "debian-faq",
	}
	if err := s.Apply(); err != nil {
		t.Errorf("Apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("Test: %v", err)
	}
}
