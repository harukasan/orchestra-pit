package homebrew_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/state/packagemanager/homebrew"
)

func TestHomebrew(t *testing.T) {
	if err := homebrew.Tap("orchestra-pit/fake"); err != nil {
		t.Errorf("Tap: %v", err)
	}

	if err := homebrew.Install("orchestra-pit/fake/stub", nil); err != nil {
		t.Errorf("Install: %v", err)
	}

	if err := homebrew.IsInstalled("orchestra-pit/fake/stub", "", nil); err != nil {
		t.Errorf("IsInstalled: %v", err)
	}

	if err := homebrew.IsInstalled("orchestra-pit/fake/stub", "0.3", nil); err != nil {
		t.Errorf("IsInstalled: %v", err)
	}

	if err := homebrew.Uninstall("orchestra-pit/fake/stub"); err != nil {
		t.Errorf("Uninstall: %v", err)
	}

	if err := homebrew.IsNotInstalled("orchestra-pit/fake/stub"); err != nil {
		t.Errorf("IsNotInstalled: %v", err)
	}

	if err := homebrew.Install("orchestra-pit/fake/stub", []string{"--with-something-awesome"}); err != nil {
		t.Errorf("Install: %v", err)
	}

	if err := homebrew.IsInstalled("orchestra-pit/fake/stub", "0.3", []string{"--with-something-awesome"}); err != nil {
		t.Errorf("IsInstalled: %v", err)
	}

	if err := homebrew.Uninstall("orchestra-pit/fake/stub"); err != nil {
		t.Errorf("Uninstall: %v", err)
	}

	if err := homebrew.Untap("orchestra-pit/fake"); err != nil {
		t.Errorf("Untap: %v", err)
	}
}
