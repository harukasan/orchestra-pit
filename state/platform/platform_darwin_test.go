// +build darwin

package platform_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/state/platform"
)

func TestDetect(t *testing.T) {
	f, err := platform.Identify()
	if err != nil {
		t.Errorf("got error, %v", err)
	}
	t.Log(f)
}
