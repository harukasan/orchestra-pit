// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build darwin

/*
Package homebrew provides command interface of brew command of Homebrew.
*/
package homebrew

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// Package describes information of the package.
type Package struct {
	Name      string             `json:"name"`
	FullName  string             `json:"full_name"`
	Desc      string             `json:"desc"`
	Homepage  string             `json:"homepage"`
	Versions  PackageVersions    `json:"versions"`
	Revision  int                `json:"revision"`
	Installed []InstalledVersion `json:"installed"`
	LinkedKeg string             `json:"linked_keg"`
	Options   []PackageOption    `json:"options"`
}

// PackageVersions describes the available versions of the package.
type PackageVersions struct {
	Stable string `json:"stable"`
	Bottle bool   `json:"bottle"`
	Devel  string `json:"devel"`
	Head   string `json:"head"`
}

// InstalledVersion describes the installed versions of the package.
type InstalledVersion struct {
	Version     string   `json:"version"`
	UsedOptions []string `json:"used_options"`
}

// PackageOption describes information of the option of the package
type PackageOption struct {
	Option      string `json:"option"`
	Description string `json:"description"`
}

// Path specifies the file path of Homebrew binary.
var Path = "/usr/local/bin/brew"

// Info returns information of the named package.
func Info(name string) (*Package, error) {
	out, err := exec.Command(Path, "info", "--json=v1", name).CombinedOutput()
	if err != nil {
		return nil, err
	}
	var packages []Package
	if err := json.Unmarshal(out, &packages); err != nil {
		return nil, err
	}
	if len(packages) > 0 {
		return &packages[0], nil
	}
	return &Package{}, notFound(fmt.Sprintf("package %s is not found", name))
}

// IsExist tests whether the named package exists.
func IsExist(name string) (bool, error) {
	_, err := Info(name)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Install executes the command to install the named package with the given
// options.
func Install(name string, options []string) error {
	options = append([]string{"install", name}, options...)
	cmd := exec.Command(Path, options...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Uninstall executes the command to uninstall the named package.
func Uninstall(name string) error {
	cmd := exec.Command(Path, "uninstall", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// IsInstalled tests whether the named package is installed on the system.
// If the named package is not installed, or the execution fails, it returns
// an error.
//
// If the version is specified, it checks whether the specified version is
// linked.
//
// If the options are specified, it also checks whether the linked version was
// built with given options.
//
// To test whether the package is NOT installed, Use IsNotInstalled function
// instead of this.
func IsInstalled(name string, version string, options []string) error {
	pkg, err := Info(name)
	if err != nil {
		if isNotFound(err) {
			return errors.New("the package is not found")
		}
		return err
	}

	// check whether the package is installed
	if len(pkg.Installed) == 0 {
		return errors.New("the package is not installed")
	}

	// check linked version
	if version != "" && version != pkg.LinkedKeg {
		return fmt.Errorf("expected version %s, but %s is installed", version, pkg.LinkedKeg)
	}

	// check build options of the linked version
	if len(options) > 0 {
		// retrieve linked version
		var linked InstalledVersion
		for _, v := range pkg.Installed {
			if v.Version == pkg.LinkedKeg {
				linked = v
				break
			}
		}

		// check whether the linked version was build with options
		for _, opt := range options {
			used := false
			for _, o := range linked.UsedOptions {
				if o == opt {
					used = true
					break
				}
			}
			if !used {
				return fmt.Errorf("the option %s is not used", opt)
			}
		}
	}

	return nil
}

// IsNotInstalled tests whether the named package is not installed on the system.
// If the package is installed or execution fails, it returns an error.
func IsNotInstalled(name string) error {
	pkg, err := Info(name)
	if err != nil {
		if isNotFound(err) {
			return nil
		}
		return err
	}
	if len(pkg.Installed) == 0 {
		return nil
	}
	return errors.New("the package is installed")
}

// Tap executes the tap command with given repository name, if the given named
// repository is not tapped.
func Tap(name string) error {
	cmd := exec.Command(Path, "tap", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Untap executes the untap command with given repository name, if the given
// named repository is tapped.
func Untap(name string) error {
	cmd := exec.Command(Path, "untap", name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// IsTapped tests whether the given named repository is tapped.
func IsTapped(name string) (bool, error) {
	out, err := exec.Command(Path, "tap").Output()
	if err != nil {
		return false, err
	}
	if bytes.Contains(out, []byte(name)) {
		return true, nil
	}
	return false, nil
}

// Update executes the update of Homebrew.
func Update() error {
	return exec.Command(Path, "update").Run()
}
