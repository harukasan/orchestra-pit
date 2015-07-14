package resource

import (
	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/opit/resource/file"
	"github.com/harukasan/orchestra-pit/opit/resource/packagemanager"
)

type Resource interface {
	States() ([]commands.State, error)
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
