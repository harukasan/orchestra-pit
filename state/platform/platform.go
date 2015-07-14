// Copyright 2015 MICHII Shunsuke. All rights reserved.

/*
Package platform implements gathering information about the platform.
*/
package platform

// Name represents the platform name.
type Name string

// Platform definitions
const (
	PlatformOSX       Name = "mac_os_x"
	PlatformDebian    Name = "debian"
	PlatformUbuntu    Name = "ubuntu"
	PlatformLinuxMint Name = "linuxmint"
	PlatformCentOS    Name = "centos"
	PlatformRHEL      Name = "rhel"
)

// Family represents the family of platforms.
type Family string

// Platform family definitions
const (
	FamilyUnknown Family = "unknown"  // unknown family
	FamilyDebian  Family = "debian"   // debian, ubuntu, linuxmint
	FamilyOSX     Family = "mac_os_x" // only osx
	FamilyRHEL    Family = "rhel"     // centos
)

// Info represents the facts about the platform.
type Info struct {
	Platform     Name
	Family       Family
	Version      string
	BuildVersion string
}

// Get returns the named attribute of the platform information as a string.
func (i *Info) Get(name string) string {
	switch name {
	case "platform":
		return string(i.Platform)
	case "family":
		return string(i.Family)
	case "version":
		return i.Version
	case "build_version":
		return i.BuildVersion
	}
	return ""
}

// LSBInfo represents the platform information of which supports LSB.
type LSBInfo struct {
	ID          string
	Release     string
	Codename    string
	Description string
}

// Get returns the named attribute of the platform information as a string.
func (i *LSBInfo) Get(name string) string {
	switch name {
	case "id":
		return i.ID
	case "release":
		return i.Release
	case "codename":
		return i.Codename
	case "Description":
		return i.Description
	}
	return ""
}
