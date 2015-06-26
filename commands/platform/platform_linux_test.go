// +build linux

package platform_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/commands/platform"
)

func TestIdentify(t *testing.T) {
	info, err := platform.Identify()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	t.Log(info)
}
