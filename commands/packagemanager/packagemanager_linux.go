// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package packagemanager

import (
	"errors"

	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/commands/platform"
)

func (s *Installed) stateForSpecificPlatform() (commands.State, error) {
	p, err := platform.Identify()
	if err != nil {
		return nil, err
	}
	switch platform.Family(p.Get("family")) {
	case platform.FamilyDebian:
		return &InstalledForDebian{s}, nil
	}
	return nil, errors.New("unsupported platform")
}

func (s *Installed) apply() error {
	ps, err := s.stateForSpecificPlatform()
	if err != nil {
		return err
	}
	return ps.Apply()
}

func (s *Installed) test() error {
	ps, err := s.stateForSpecificPlatform()
	if err != nil {
		return err
	}
	return ps.Test()
}

func (s *Removed) stateForSpecificPlatform() (commands.State, error) {
	p, err := platform.Identify()
	if err != nil {
		return nil, err
	}
	switch platform.Family(p.Get("family")) {
	case platform.FamilyDebian:
		return &RemovedForDebian{s}, nil
	}
	return nil, errors.New("unsupported platform")
}

func (s *Removed) apply() error {
	ps, err := s.stateForSpecificPlatform()
	if err != nil {
		return err
	}
	return ps.Apply()
}

func (s *Removed) test() error {
	ps, err := s.stateForSpecificPlatform()
	if err != nil {
		return err
	}
	return ps.Test()
}
