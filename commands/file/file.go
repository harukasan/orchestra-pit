//
// ## States
//
//
// ## TODO
//
// - Template
//

/*
Package file implements the commands to manage the state of file.

Following states are implemented:

	- Copy ... manages the file whose contents is a copy of the source file
  - Directory ... manages the directory existence
  - Hardlink ... manages the hard link file
  - Symlink ... manages the symbolic link file
  - Owner ... manages owner and group of the file
	- Mode ... manages the file mode and permissions.

*/
package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/harukasan/orchestra-pit/commands"
)

// Copy manages the file whose content is a copy of the src file.
//
// Name specifies the requesting file name. Copy keeps content of the file to
// copy of content of the src file.
//
// Src specifies the requesting original source file.
//
// If the Backup value is not empty, create the backup file to get the original
// file back.
type Copy struct {
	Name   string
	Src    string
	Backup string
}

// NewCopyState returns a new Copy state with the given options. If the given
// options is not valid, NewCopyState returns nil and an error.
//
// NewCopyState has following options:
//
//   - name ... specifies the name of the target file
//   - src ... specifies the original source file
//
func NewCopyState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	src := options.Get("src")
	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	if src == "" {
		return nil, fmt.Errorf("src: not specified")
	}

	return &Copy{
		Name: name,
		Src:  src,
	}, nil
}

// Apply tries to copy the content from the src file. When the backup value is
// not empty, rename the file to given the backup name before the copying file.
func (s *Copy) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)

	r, err := os.Open(s.Src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.Rename(s.Name, s.Backup); err != nil {
		if os.IsExist(err) {
			return err
		}
	}

	w, err := os.Create(s.Name)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

// Test tests whether the file contains the same contents of the src file.
func (s *Copy) Test() error {
	dest, err := os.Open(s.Name)
	if err != nil {
		return err
	}
	defer dest.Close()

	src, err := os.Open(s.Src)
	if err != nil {
		return err
	}
	defer src.Close()

	equal, err := compareFile(src, dest)
	if err != nil {
		return err
	}
	if !equal {
		return errors.New("content of the file is different to the src file")

	}

	return nil
}

func compareFile(a, b io.Reader) (bool, error) {
	buf := make([]byte, 64*1024)
	bufa, bufb := buf[0:32*1024], buf[32*1024:]

	for {
		na, aerr := a.Read(bufa)
		nb, berr := b.Read(bufb)
		if na != nb {
			return false, nil
		}
		if na > 0 && !bytes.Equal(bufa, bufb) {
			return false, nil
		}
		if aerr == io.EOF && berr == io.EOF {
			break
		}
		if aerr != nil {
			return false, aerr
		}
		if berr != nil {
			return false, berr
		}
	}
	return true, nil
}

// Directory manages the directory existence.
//
// Name specifies the requesting directory name. Directory tires to keep the
// directory existence.
type Directory struct {
	Name string
}

// NewDirectoryState returns a new Directory state with the given options. If
// the given options is not valid, NewDirectoryState returns nil and an error.
//
// NewDirectoryState has following options:
//
//   - name ... specifies the name of the directory
//
func NewDirectoryState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	return &Directory{
		Name: name,
	}, nil
}

// Apply tries to make the named directory. If failed to make a directory, Apply
// returns an error.
func (s *Directory) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)

	return os.Mkdir(s.Name, 0777)
}

// Test tests whether the named file is a directory.
func (s *Directory) Test() error {
	info, err := FileInfoCache.Stat(s.Name)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return fmt.Errorf("the file is not a directory")
}

// Absence manages the file non-existence.
//
// Name the specifies the named file that should not exists.
//
type Absence struct {
	Name string
}

// NewAbsenceState returns a new Absence state of the given named file. If the
// name attribute is not specified, NewAbsenceState returns nil and an error.
//
// NewAbsenceState has following options:
//
//   - name ... specifies the name of the file that should not exists.
//
func NewAbsenceState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	return &Absence{
		Name: name,
	}, nil
}

// Apply tries to remove the file.
func (s *Absence) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)
	err := os.Remove(s.Name)
	if err != nil && os.IsExist(err) {
		return err
	}
	return nil
}

// Test tests whether the named file does not exists.
func (s *Absence) Test() error {
	_, err := FileInfoCache.Stat(s.Name)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
		return nil
	}
	return errors.New("the file exists")
}
