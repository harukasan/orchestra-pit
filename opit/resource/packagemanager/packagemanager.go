/*
Package packagemanager implements the applying state of packages.
*/
package packagemanager

import (
	"fmt"

	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/state"
	"github.com/harukasan/orchestra-pit/state/packagemanager"
)

// Resource represents the attributes of package resource.
type Resource struct {
	Desc    string   `json:"desc" yaml:"desc"`
	Name    string   `json:"name" yaml:"name"`
	Version string   `json:"version" yaml:"version"`
	Options []string `json:"options" yaml:"options"`
	Update  bool     `json:"update" yaml:"update"`
	State   string   `json:"state" yaml:"state"`
}

func (r *Resource) States() ([]state.State, error) {
	states := []state.State{}

	if r.State == "" {
		r.State = "installed"
		logger.Debugf(`parameter "state" is not specified, assume as "%s"`, r.State)
	}

	var s state.State
	var err error
	switch r.State {
	case "installed":
		s, err = r.installedState()
	case "removed":
		s, err = r.removedState()
	}
	if err != nil {
		return nil, err
	}
	states = append(states, s)

	return states, nil
}

func (r *Resource) installedState() (state.State, error) {
	if r.Name == "" {
		return nil, fmt.Errorf(`parameter "name" is required`)
	}
	return &packagemanager.Installed{
		Name:    r.Name,
		Version: r.Version,
		Options: r.Options,
		Update:  r.Update,
	}, nil
}

func (r *Resource) removedState() (state.State, error) {
	if r.Name == "" {
		return nil, fmt.Errorf(`parameter "name" is required`)
	}
	return &packagemanager.Installed{
		Name: r.Name,
	}, nil
}
