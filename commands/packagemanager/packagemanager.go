// Copyright 2015 MICHII Shunsuke. All rights reserved.

/*
Package packagemanager implements the commands to manage the packages of the
package management system (PMS).

Following PMSes is supported now:

  - Homebrew (Mac OS X)
	- APT      (Debian and its derivatives)

*/
package packagemanager

// Installed tries to keep that the named package is installed on the system.
//
// Name and Version specifies the name and version of package. If Version is not
// specified, the latest version should be installed.
//
// Options is passed as arguments of the command of package management system.
//
// If Update flag is set true, Apply tries to update the package to the
// latest version.
type Installed struct {
	Name    string
	Version string
	Options []string
	Update  bool
}

// Apply tries to install the named package into the system. If failed to
// install the package, it returns an error.
func (s *Installed) Apply() error {
	return s.apply()
}

// Test tests whether the package is installed. If the package is not installed
// or the another version is installed, it returns an error.
func (s *Installed) Test() error {
	return s.test()
}

// Removed tries to keep that the named package is removed on the system.
//
// Name specifies the name of package.
type Removed struct {
	Name string
}

// Apply tries to remove the named package from the system. If failed to
// remiving the package, it returns an error.
func (s *Removed) Apply() error {
	return s.apply()
}

// Test tests whether the package is not installed. If the named package is not
// absent, it returns an error.
func (s *Removed) Test() error {
	return s.test()
}
