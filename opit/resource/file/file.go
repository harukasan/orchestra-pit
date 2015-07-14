/*
Package file implements the applying state of file resources.
*/
package file

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/commands/file"
	"github.com/harukasan/orchestra-pit/opit/logger"
)

type Resource struct {
	Desc   string `json:"desc" yaml:"desc"`
	Path   string `json:"path"  yaml:"path"`
	State  string `json:"state" yaml:"state"`
	Src    string `json:"src"   yaml:"src"`
	Backup string `json:"backup" yaml:"backup"`
	Mode   string `json:"mode"  yaml:"mode"`
}

func (r *Resource) States() ([]commands.State, error) {
	states := []commands.State{}

	if r.State == "" {
		r.State = "file"
		logger.Debugf(`parameter "state" is not specified, assume as "%s"`, r.State)
	}
	if stateFuncMap[r.State] != nil {
		s, err := stateFuncMap[r.State](r)
		if err != nil {
			return nil, err
		}
		states = append(states, s)
	}

	if r.Mode != "" {
		s := &file.Mode{
			Name: r.Path,
			Mode: r.Mode,
		}
		states = append(states, s)
	}

	return states, nil
}

type stateFunc func(r *Resource) (commands.State, error)

var stateFuncMap = map[string]stateFunc{
	"absence":   absenceState,
	"directory": directoryState,
	"file":      fileState,
	"hardlink":  hardlinkState,
	"symlink":   symlinkState,
}

func absenceState(r *Resource) (commands.State, error) {
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}
	return &file.Absence{
		Name: r.Path,
	}, nil
}

func directoryState(r *Resource) (commands.State, error) {
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}
	return &file.Directory{
		Name: r.Path,
	}, nil
}

func fileState(r *Resource) (commands.State, error) {
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}
	if r.Src == "" {
		if strings.HasPrefix(r.Path, "/") {
			wd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			r.Src = path.Join(wd, "files", r.Path[1:])
		}
		logger.Debugf(`parameter "src" is not specified, assume as "%s"`, r.Src)
	}
	if r.Backup != "" {
		if !strings.ContainsRune(r.Backup, '/') {
			r.Backup = path.Join(path.Dir(r.Path), r.Backup)
		}
	}
	return &file.Copy{
		Name:   r.Path,
		Src:    r.Src,
		Backup: r.Backup,
	}, nil
}

func hardlinkState(r *Resource) (commands.State, error) {
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}
	if r.Src == "" {
		return nil, fmt.Errorf(`parameter "src" is required`)
	}
	return &file.Hardlink{
		Name: r.Path,
		Src:  r.Src,
	}, nil
}

func symlinkState(r *Resource) (commands.State, error) {
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}
	if r.Src == "" {
		return nil, fmt.Errorf(`parameter "src" is required`)
	}
	return &file.Symlink{
		Name: r.Path,
		Src:  r.Src,
	}, nil
}
