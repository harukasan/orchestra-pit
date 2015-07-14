// +build linux darwin dragonfly freebsd openbsd netbsd solaris

package file_test

import (
	"os"
	"os/user"
	"strconv"
	"testing"

	"github.com/harukasan/orchestra-pit/state/file"
)

type options map[string]string

func (o options) Get(name string) string {
	return o[name]
}

func TestHardlink(t *testing.T) {
	src := d.MakeDummyFile("source_file")

	s := &file.Hardlink{
		Name: d.NewFilePath("dest"),
		Src:  src,
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}

func TestHardlinkTestShouldFailWhenTheTargetFileIsNotFound(t *testing.T) {
	s := &file.Hardlink{
		Name: d.NewFilePath("the_target_file_is_not_found"),
		Src:  d.NewFilePath("the_source_file_is_not_found"),
	}
	err := s.Test()
	if err == nil {
		t.Errorf("got no error when the file is not found")
	}
	t.Log(err)
}

func TestHardlinkTestShouldFailWhenTheSourceFileIsNotFound(t *testing.T) {
	dest := d.MakeDummyFile("the_destination_file")
	s := &file.Hardlink{
		Name: dest,
		Src:  d.NewFilePath("the_source_file_is_not_found"),
	}
	err := s.Test()
	if err == nil {
		t.Errorf("got no error when the file points to nothing file")
	}
	t.Log(err)
}

func TestHardlinkTestShouldFailWhenTheFilePointsToTheAnotherFile(t *testing.T) {
	src1 := d.MakeDummyFile("the_expect_source_file")
	src2 := d.MakeDummyFile("the_another_source_file")
	dest := d.NewFilePath("the_subject")
	os.Link(src2, dest)

	s := &file.Hardlink{
		Name: dest,
		Src:  src1,
	}
	err := s.Test()
	if err == nil {
		t.Errorf("got no error when the file points to the another file")
	}
	t.Log(err)
}

func TestSymlink(t *testing.T) {
	src := d.MakeDummyFile("test_symlink_")

	s := &file.Symlink{
		Name: d.NewFilePath("test_symlink_target"),
		Src:  src,
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}

func TestSymLinkTestShouldFailWhenTheFileIsARegularFile(t *testing.T) {
	dest := d.MakeDummyFile("the_regular_file")
	src := d.MakeDummyFile("the_source_file")

	s := &file.Symlink{
		Name: dest,
		Src:  src,
	}
	err := s.Test()
	if err == nil {
		t.Errorf("got no error when the file is not a symbolic link")
	}
	t.Log(err)
}

func TestSymLinkTestShouldFailWhenTheFilePointsToTheAnotherFile(t *testing.T) {
	src1 := d.MakeDummyFile("the_source_file")
	src2 := d.MakeDummyFile("the_another_file")
	dest := d.NewFilePath("the_symbolic_link_file")

	os.Symlink(src2, dest)

	s := &file.Symlink{
		Name: dest,
		Src:  src1,
	}
	err := s.Test()
	if err == nil {
		t.Errorf("got no error when the file points to the another file")
	}
	t.Log(err)
}

func TestOwner(t *testing.T) {
	target := d.MakeDummyFile("test_owner_")

	if os.Getuid() != 0 {
		t.Skip("the test about the owner state is skipped, required running as a root")
	}

	nobody, _ := user.Lookup("nobody")
	nobodyUid, _ := strconv.Atoi(nobody.Uid)
	nobodyGid, _ := strconv.Atoi(nobody.Gid)

	s := &file.Owner{
		Name: target,
		Uid:  uint32(nobodyUid),
		Gid:  uint32(nobodyGid),
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}

func TestMode(t *testing.T) {
	target := d.MakeDummyFile("test_mode_")

	s := &file.Mode{
		Name: target,
		Mode: "600",
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}
