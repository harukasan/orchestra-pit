package file_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/harukasan/orchestra-pit/state/file"
	"github.com/harukasan/orchestra-pit/state/file/testutil"
)

var d = testutil.TempDir()

func TestCopy(t *testing.T) {
	src := d.MakeDummyFile("test_copy_")
	dest := d.NewFilePath("test_copy_target")

	s := &file.Copy{
		Name:   dest,
		Src:    src,
		Backup: "",
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}

	// with backup
	src2 := d.MakeDummyFile("test_copy_")
	backup := d.NewFilePath("test_copy_backup")
	s = &file.Copy{
		Name:   dest,
		Src:    src2,
		Backup: backup,
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}

	bb, err := ioutil.ReadFile(backup)
	if err != nil {
		t.Errorf("failed to read backed up file, %v", err)
	}
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		t.Errorf("failed to read src file, %v", err)
	}
	if !bytes.Equal(bb, bs) {
		t.Errorf("the backed up file has different content to the src file")
	}
}

func TestDirectory(t *testing.T) {
	s := &file.Directory{
		Name: d.NewFilePath("test_directory_target"),
	}

	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}

func TestAbsence(t *testing.T) {
	s := &file.Absence{
		Name: d.NewFilePath("non_existence_file"),
	}
	if err := s.Apply(); err != nil {
		t.Errorf("got error on apply: %v", err)
	}
	if err := s.Test(); err != nil {
		t.Errorf("got error on test: %v", err)
	}
}
