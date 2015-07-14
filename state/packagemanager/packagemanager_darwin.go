// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build darwin

package packagemanager

import (
	"strings"
	"sync"

	"github.com/harukasan/orchestra-pit/state/packagemanager/homebrew"
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
		if err := homebrew.Update(); err != nil {
			return err
		}
		u.updated = true
		return nil
	}
	u.RUnlock()
	return nil
}

var update = &updateOnce{}

func (s *Installed) apply() error {
	if s.Update {
		if err := update.do(); err != nil {
			return err
		}
	}

	// if the package name contains the slash, tap the repository.
	if i := strings.LastIndex(s.Name, "/"); i > 0 {
		if err := homebrew.Tap(s.Name[0:i]); err != nil {
			return err
		}
	}

	// try to uninstall the named package if the different version is installed.
	if err := homebrew.IsInstalled(s.Name, s.Version, s.Options); err == nil {
		err := homebrew.Uninstall(s.Name)
		if err != nil {
			return err
		}
	}

	return homebrew.Install(s.Name, s.Options)
}

func (s *Installed) test() error {
	return homebrew.IsInstalled(s.Name, s.Version, s.Options)
}

func (s *Removed) apply() error {
	return homebrew.Uninstall(s.Name)
}

func (s *Removed) test() error {
	if err := homebrew.IsNotInstalled(s.Name); err != nil {
		return err
	}
	return nil
}
