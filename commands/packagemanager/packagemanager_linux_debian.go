// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package packagemanager

import (
	"sync"

	"github.com/harukasan/orchestra-pit/commands/packagemanager/apt"
)

type updateOnce struct {
	sync.RWMutex
	updated bool
}

func (u *updateOnce) do() error {
	u.RLock()
	if !u.updated {
		u.RUnlock()
		u.Lock()
		defer u.Unlock()
		if err := apt.Update(); err != nil {
			return err
		}
		u.updated = true
		return nil
	}
	u.RUnlock()
	return nil
}

var update = &updateOnce{}

// InstalledForDebian implements state of which the package is installed for
// the platform of Debian or its derivatives.
type InstalledForDebian struct {
	*Installed
}

// Apply tries to install the package with APT for Debian or its derivatives. If
// the installation fails, it returns an error.
func (s *InstalledForDebian) Apply() error {
	if s.Update {
		if err := update.do(); err != nil {
			return err
		}
	}
	return apt.Install(s.Name, s.Version)
}

// Test checks whether the package is successfully installed on the platform of
// Debian or its derivatives.
func (s *InstalledForDebian) Test() error {
	return apt.IsInstalled(s.Name, s.Version)
}

// RemovedForDebian implements state of which the package is removed for the
// platform of Debian or its derivatives.
type RemovedForDebian struct {
	*Removed
}

// Apply tries to remove the package with APT on the Debian or its derivatives.
// If the removing the package fails, it returns an error.
func (s *RemovedForDebian) Apply() error {
	return apt.Remove(s.Name)
}

// Test checks whether the package is absent on the Debian or its derivatives.
func (s *RemovedForDebian) Test() error {
	return apt.IsNotInstalled(s.Name)
}
