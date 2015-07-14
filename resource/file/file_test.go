package file_test

import (
	"strings"
	"testing"

	"github.com/harukasan/orchestra-pit/resource/file"
	filestate "github.com/harukasan/orchestra-pit/state/file"
)

func TestStatesWithMinimumArguments(t *testing.T) {
	r := &file.Resource{
		Path: "/tmp/test",
	}

	states, err := r.States()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if got := len(states); got != 1 {
		t.Errorf("got %d states, exected just 1", got)
	}

	if s, ok := states[0].(*filestate.Copy); ok {
		if s.Name != r.Path {
			t.Errorf("got Name %v, expected %v", s.Name, r.Path)
		}
		if !strings.HasSuffix(s.Src, r.Path[1:]) {
			t.Errorf("got Src %v, expected has following suffix: %v", s.Name, r.Path[1:])
		}
	} else {
		t.Errorf("state is not a Copy state")
	}
}

func TestCopyState(t *testing.T) {
	r := &file.Resource{
		Path:   "/tmp/test",
		State:  "file",
		Src:    "/tmp/src",
		Backup: "/tmp/backup",
	}

	states, err := r.States()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if got := len(states); got != 1 {
		t.Errorf("got %d states, exected just 1", got)
	}

	if s, ok := states[0].(*filestate.Copy); ok {
		if s.Name != r.Path {
			t.Errorf("got Name %v, expected %v", s.Name, r.Path)
		}
		if s.Src != r.Src {
			t.Errorf("got Src %v, expected %v", s.Src, r.Src)
		}
		if s.Backup != r.Backup {
			t.Errorf("got Backup %v, expected %v", s.Backup, r.Backup)
		}
	} else {
		t.Errorf("state is not a Copy state")
	}
}
