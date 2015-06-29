// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package apt

import (
	"bytes"
	"errors"
	"os"

	"github.com/harukasan/orchestra-pit/commands/exec"
)

// APTGetPath specifies the file path of the apt-get command
var APTGetPath = "/usr/bin/apt-get"

// DPKGQueryPath specifies the file path of the dpkg-query command
var DPKGQueryPath = "/usr/bin/dpkg-query"

// InstallOptions is array of argument which passed to apt-get install command.
var InstallOptions = []string{
	"-o Dpkg::Options::='--force-confdef'",
	"-o Dpkg::Options::='--force-confold'",
}

// Install executes the apt-get command to install the named package. If the
// version is specified, the version is also passed to apt-get command.
func Install(name string, version string) error {
	args := []string{"install", "-y"}
	args = append(args, InstallOptions...)
	if version != "" {
		args = append(args, exec.ShellEscape(name+"="+version))
	} else {
		args = append(args, exec.ShellEscape(name))
	}
	cmd := exec.Command(APTGetPath, args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	return cmd.Run()
}

// Remove executes the apt-get command to remove the named package.
func Remove(name string) error {
	cmd := exec.Command(APTGetPath, "remove", "-y", exec.ShellEscape(name))
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	return cmd.Run()
}

// IsInstalled tests whether the named package is installed on the system.
// If the named package is not installed, or the execution fails, it returns
// an error.
//
// If the version is specified, it checks whether the specified version is
// linked.
//
// To test whether the package is NOT installed, Use IsNotInstalled function
// instead of this.
func IsInstalled(name string, version string) error {
	cmd := exec.Command(DPKGQueryPath, "--showformat=${Status}\\n${Version}", "--show", exec.ShellEscape(name))
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if !bytes.HasPrefix(out, []byte("install ok installed\n")) {
		return errors.New("the package is not installed")
	}
	if version != "" {
		if i := bytes.IndexRune(out, '\n'); i > 0 {
			if !bytes.HasPrefix(out[i+1:], []byte(version)) {
				return errors.New("the different version is installed")
			}
		} else {
			return errors.New("failed to parse the result of dpkg-query")
		}
	}
	return nil
}

// IsNotInstalled tests whether the named package is not installed on the system.
// If the package is installed or execution fails, it returns an error.
func IsNotInstalled(name string) error {
	cmd := exec.Command(DPKGQueryPath, "--show", exec.ShellEscape(name))
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// TODO:
		// I could not find the apt/dpkg command to retrieve status of uninstalled
		// package. just check output message.
		if bytes.Contains(out, []byte("no packages found")) {
			return nil
		}
		return err
	}
	return errors.New("the package is installed")
}

// Update executes updating lists of apt packages.
func Update() error {
	return exec.Command(APTGetPath, "update").Run()
}
