package resource

import (
	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/resource/file"
	"github.com/harukasan/orchestra-pit/resource/packagemanager"
	"github.com/harukasan/orchestra-pit/state"
)

type Resource interface {
	States() ([]state.State, error)
}

func New(t string) Resource {
	switch t {
	case "file":
		return &file.Resource{}
	case "package":
		return &packagemanager.Resource{}
	}
	return nil
}

func Apply(r Resource) error {
	states, err := r.States()
	if err != nil {
		return err
	}
	for _, state := range states {
		logger.Debugf("applying state: %s", state)
		if err := state.Apply(); err != nil {
			return err
		}
	}
	return nil
}

func Test(r Resource) error {
	states, err := r.States()
	if err != nil {
		return err
	}
	for _, state := range states {
		logger.Debugf("testing state: %s", state)
		if err := state.Test(); err != nil {
			return err
		}
	}
	return nil
}
