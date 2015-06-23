// +build linux darwin dragonfly freebsd openbsd netbsd solaris

package file

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/harukasan/orchestra-pit/commands"
)

// Hardlink manages the hard link existence and where the file points to.
//
// Name specifies the file name. Hardlink creates the hard link file with
// the specified name.
//
// Src specifies the source file of the hard link. The link file points to the
// directory entry of the source file.
type Hardlink struct {
	Name string
	Src  string
}

// NewHardlinkState returns a new Hardlink state with the given options. If the
// given options is not valid, NewHardlinkState returns nil and an error.
//
// NewHardlinkState has following options:
//
//   - name ... specifies the name of the hard link file
//   - src ... specifies the source file which is pointed from the hard link
//
func NewHardlinkState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	src := options.Get("src")
	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	if src == "" {
		return nil, fmt.Errorf("src: not specified")
	}

	return &Hardlink{
		Name: name,
		Src:  src,
	}, nil
}

// Apply tries to make a hardlink to src file.
func (s *Hardlink) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)
	return os.Link(s.Src, s.Name)
}

// Test tests whether the file points to same location as the src file.
func (s *Hardlink) Test() error {
	destInfo, err := FileInfoCache.Stat(s.Name)
	if err != nil {
		return err
	}
	srcInfo, err := FileInfoCache.Stat(s.Src)
	if err != nil {
		return err
	}
	if destStat, ok := destInfo.Sys().(*syscall.Stat_t); ok {
		if srcStat, ok := srcInfo.Sys().(*syscall.Stat_t); ok {
			if destStat.Ino != srcStat.Ino {
				return fmt.Errorf("the file is point to the another inode")
			}
			return nil
		}
	}
	return fmt.Errorf("failed to get file stat")
}

// Symlink manages the named symbolick link file and where the file links to.
//
// Name specifies the file name. Symlink creates the symbolic link file with
// specified the name.
//
// Src speicfies the source file of the symbolic link. The link file points to
// the source file name.
type Symlink struct {
	Name string
	Src  string
}

// NewSymlinkState returns a new Symlink state with the given options. If the
// given options is not valid, NewSymlinkState returns nil and an error.
//
// NewSymlinkState has following options:
//
//   - name ... specifies the name of the symbolic link file
//   - src ... specifies the requesting source file which is pointed from the link
//
func NewSymlinkState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	src := options.Get("src")
	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	if src == "" {
		return nil, fmt.Errorf("src: not specified")
	}

	return &Symlink{
		Name: name,
		Src:  src,
	}, nil
}

// Apply tries to make a symbolic link which is linked to the Src.
func (s *Symlink) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)

	return os.Symlink(s.Src, s.Name)
}

// Test tests whether the file points to the Src.
func (s *Symlink) Test() error {
	fact, err := os.Readlink(s.Name)
	if err != nil {
		return err
	}
	if fact != s.Src {
		return fmt.Errorf("the file is not linked to %s, but %s", s.Src, fact)
	}
	return nil
}

// Owner manages owner and group of the named file.
//
// Name speicifes the file name.
//
// Uid and Gid specifies the ID of owner and group of the file. To specify the
// owner and group as string, use NamedOwner insteads of Owner state.
//
type Owner struct {
	Name string
	Uid  uint32
	Gid  uint32
}

// NewOwnerState returns a new Owner state with the given options. If the
// given options is not valid, NewOwnerState returns nil and an error.
//
// NewOwnerState have below options:
//
//   - name ... specifies the name of target file
//   - uid ... specifies the requesting uid of the file owner
//   - gid ... specifies the requesting gid of the group that owns the file
//
func NewOwnerState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	uidString := options.Get("uid")
	gidString := options.Get("gid")

	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	if uidString == "" {
		return nil, fmt.Errorf("uid: not specified")
	}
	if gidString == "" {
		return nil, fmt.Errorf("gid: not specified")
	}

	uid, err := strconv.Atoi(uidString)
	if err != nil {
		return nil, fmt.Errorf("uid: %s", err)
	}

	gid, err := strconv.Atoi(gidString)
	if err != nil {
		return nil, fmt.Errorf("uid: %s", err)
	}

	return &Owner{
		Name: name,
		Uid:  uint32(uid),
		Gid:  uint32(gid),
	}, nil
}

// Apply tries to change the file owner and group.
func (s *Owner) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)

	return os.Chown(s.Name, int(s.Uid), int(s.Gid))
}

// Test tests whether the owner and group of the is requested.
func (s *Owner) Test() error {
	info, err := FileInfoCache.Stat(s.Name)
	if err != nil {
		return fmt.Errorf("faild to get stat on testing file owner: %v", err)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if stat.Uid != s.Uid {
			return fmt.Errorf("wrong uid, requested: %d, but %d", s.Uid, stat.Uid)
		}
		if stat.Gid != s.Gid {
			return fmt.Errorf("wrong gid, requested: %d, but %d", s.Gid, stat.Gid)
		}
	}
	return nil
}

// Mode manages the file mode and permissions.
//
// Name specifies the name of the file.
//
// Mode specifies the file mode of the file.
type Mode struct {
	Name string
	Mode string
}

// NewModeState returns a new Mode state with the given options. If the given
// options is not valid, NewModeState returns nil and an error.
//
// NewModeState has the following options:
//
//   - name ... specifies the name of the file
//   - mode ... specifies the requesting mode of the file.
//
func NewModeState(options commands.Options) (commands.State, error) {
	name := options.Get("name")
	mode := options.Get("mode")

	if name == "" {
		return nil, fmt.Errorf("name: not specified")
	}
	if mode == "" {
		return nil, fmt.Errorf("mode: not specified")
	}
	if _, err := ParseMode(mode, os.FileMode(0)); err != nil {
		return nil, err
	}

	return &Mode{
		Name: name,
		Mode: mode,
	}, nil
}

// Apply tries to keep the file mode to the requested mode.
func (s *Mode) Apply() error {
	FileInfoCache.Lock()
	defer FileInfoCache.ClearAndUnlock(s.Name)

	fi, err := os.Stat(s.Name)
	if err != nil {
		return err
	}
	mode, err := ParseMode(s.Mode, fi.Mode())
	if err != nil {
		return err
	}
	return os.Chmod(s.Name, mode)
}

// Test tests whether the file mode is requested.
func (s *Mode) Test() error {
	fi, err := FileInfoCache.Stat(s.Name)
	if err != nil {
		return err
	}
	mode, err := ParseMode(s.Mode, fi.Mode())
	if mode != fi.Mode() {
		return fmt.Errorf("the mode %s is different from the requested: %s", fi.Mode(), s.Mode)
	}
	return nil
}
