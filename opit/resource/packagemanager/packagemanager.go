/*
Package packagemanager implements the applying state of packages.
*/
package packagemanager

import (
	"fmt"

	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/commands/packagemanager"
	"github.com/harukasan/orchestra-pit/opit/logger"
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

func (r *Resource) States() ([]commands.State, error) {
	states := []commands.State{}

	if r.State == "" {
		r.State = "installed"
		logger.Debugf(`parameter "state" is not specified, assume as "%s"`, r.State)
	}

	var s commands.State
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

func (r *Resource) installedState() (commands.State, error) {
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

func (r *Resource) removedState() (commands.State, error) {
	if r.Name == "" {
		return nil, fmt.Errorf(`parameter "name" is required`)
	}
	return &packagemanager.Installed{
		Name: r.Name,
	}, nil
}
