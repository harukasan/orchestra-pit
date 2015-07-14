// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package apt_test

import (
	"testing"

	"github.com/harukasan/orchestra-pit/state/packagemanager/apt"
)

var TargetPackage = "debian-faq"

func TestAPT(t *testing.T) {
	if err := apt.Install(TargetPackage, ""); err != nil {
		t.Errorf("Install: %v", err)
	}
	if err := apt.IsInstalled(TargetPackage, ""); err != nil {
		t.Errorf("IsInstalled: %v", err)
	}
	if err := apt.Remove(TargetPackage); err != nil {
		t.Errorf("Remove: %v", err)
	}
	if err := apt.IsNotInstalled(TargetPackage); err != nil {
		t.Errorf("IsNotInstalled: %v", err)
	}
}

func TestAPTWithVersion(t *testing.T) {
	if err := apt.Install(TargetPackage, "5.0.3"); err != nil {
		t.Errorf("Install: %v", err)
	}
	if err := apt.IsInstalled(TargetPackage, ""); err != nil {
		t.Errorf("IsInstalled: %v", err)
	}
	if err := apt.Remove(TargetPackage); err != nil {
		t.Errorf("Remove: %v", err)
	}
	if err := apt.IsNotInstalled(TargetPackage); err != nil {
		t.Errorf("IsNotInstalled: %v", err)
	}
}
