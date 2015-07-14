package file_test

import (
	"os"
	"testing"

	"github.com/harukasan/orchestra-pit/state/file"
)

type testPattern struct {
	input    string
	base     os.FileMode
	expected os.FileMode
}

var patterns = []testPattern{
	testPattern{
		input:    "644",
		base:     os.FileMode(0777),
		expected: os.FileMode(0644),
	},
	testPattern{
		input:    "go-w",
		base:     os.FileMode(0777),
		expected: os.FileMode(0755),
	},
	testPattern{
		input:    "=rw,+X",
		base:     os.FileMode(0000),
		expected: os.FileMode(0666),
	},
	testPattern{
		input:    "=rw,+X",
		base:     os.FileMode(0000 | os.ModeDir),
		expected: os.FileMode(0777 | os.ModeDir),
	},
	testPattern{
		input:    "u=rwx,go=rx",
		base:     os.FileMode(0000),
		expected: os.FileMode(0755),
	},
	testPattern{
		input:    "u=rwx,go=u-w",
		base:     os.FileMode(0000),
		expected: os.FileMode(0755),
	},
	testPattern{
		input:    "go=",
		base:     os.FileMode(0777),
		expected: os.FileMode(0700),
	},
	testPattern{
		input:    "g=u-w",
		base:     os.FileMode(0700),
		expected: os.FileMode(0750),
	},
	testPattern{
		input:    "o=,g=o,u=g",
		base:     os.FileMode(0753),
		expected: os.FileMode(0000),
	},
	testPattern{
		input:    "a+stx",
		base:     os.FileMode(0),
		expected: os.FileMode(0111 | os.ModeSetuid | os.ModeSetgid | os.ModeSticky),
	},
}

var invalidInputs = []string{
	"g",
	"X",
	",,,",
	"0aaa",
}

func TestParseModes(t *testing.T) {
	for _, p := range patterns {
		t.Logf("test: %s", p.input)
		got, err := file.ParseMode(p.input, p.base)
		if err != nil {
			t.Errorf("%s on %v : got error %v", p.input, p.base, err)
		}
		if p.expected != got {
			t.Errorf("%s on %v : got %v, expected %v", p.input, p.base, got, p.expected)
		}
	}
}

func TestParseModeWithInvalidParams(t *testing.T) {
	for _, i := range invalidInputs {
		t.Logf("test: %s", i)
		_, err := file.ParseMode(i, os.FileMode(0))
		if err == nil {
			t.Errorf("got no error")
		}
		t.Log(err)
	}
}
