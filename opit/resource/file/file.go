/*
Package file implements the resource state of files.
*/
package file

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/commands/file"
	"github.com/harukasan/orchestra-pit/opit/logger"
)

type Resource struct {
	base   string
	Desc   string `json:"desc" yaml:"desc"`
	Path   string `json:"path"  yaml:"path"`
	State  string `json:"state" yaml:"state"`
	Src    string `json:"src"   yaml:"src"`
	Backup string `json:"backup" yaml:"backup"`
	Mode   string `json:"mode"  yaml:"mode"`
}

func (r *Resource) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("file(")
	if r.Desc != "" {
		buf.WriteString(r.Desc)
	} else {
		buf.WriteString(fmt.Sprintf(`path="%s"`, r.Path))
		if r.State != "" {
			buf.WriteString("," + fmt.Sprintf(`state="%s"`, r.State))
		}
		if r.Src != "" {
			buf.WriteString("," + fmt.Sprintf(`src="%s"`, r.Src))
		}
		if r.Backup != "" {
			buf.WriteString("," + fmt.Sprintf(`backup="%s"`, r.Backup))
		}
		if r.Mode != "" {
			buf.WriteString("," + fmt.Sprintf(`mode="%s"`, r.Mode))
		}
	}
	buf.WriteRune(')')
	return string(buf.Bytes())
}

func (r *Resource) makeStates() ([]commands.State, error) {
	states := []commands.State{}
	if r.Path == "" {
		return nil, fmt.Errorf(`parameter "path" is required`)
	}

	if r.State == "" {
		r.State = "file"
		logger.Debugf(`parameter "state" is not specified, assume as "%s"`, r.State)
	}
	switch r.State {
	case "absence":
		s := &file.Absence{
			Name: r.Path,
		}
		states = append(states, s)

	case "directory":
		s := &file.Directory{
			Name: r.Path,
		}
		states = append(states, s)

	case "file":
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
		s := &file.Copy{
			Name:   r.Path,
			Src:    r.Src,
			Backup: r.Backup,
		}
		states = append(states, s)

	case "hardlink":
		if r.Src == "" {
			return nil, fmt.Errorf(`parameter "src" is required`)
		}
		s := &file.Hardlink{
			Name: r.Path,
			Src:  r.Src,
		}
		states = append(states, s)

	case "symlink":
		if r.Src == "" {
			return nil, fmt.Errorf(`parameter "src" is required`)
		}
		s := &file.Symlink{
			Name: r.Path,
			Src:  r.Src,
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

func (r *Resource) Apply() error {
	states, err := r.makeStates()
	if err != nil {
		return err
	}

	for _, state := range states {
		logger.Debugf("applying state: %s", state)
		err := state.Apply()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resource) Test() error {
	states, err := r.makeStates()
	if err != nil {
		return err
	}

	for _, state := range states {
		logger.Debugf("testing state: %s", state)
		err := state.Test()
		if err != nil {
			return err
		}
	}

	return nil
}
