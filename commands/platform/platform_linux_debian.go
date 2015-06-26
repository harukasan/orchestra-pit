// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package platform

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/harukasan/orchestra-pit/commands"
)

// IdentifyDebianRelease tires to read the debian_version file to
// identify that the platform is a Debian family.
// To detect the derivatives of the Debian (Ubuntu and LinuxMint),
// IdentifyDebianRelease tries to retrieve the LSB release information.
//
// If failed to identify the platform or the platform is not Debian or the
// derivatives of Debian, IdentifyDebianRelease returns an ErrNotIdentifier.
func IdentifyDebianRelease() (*Info, error) {
	file, err := os.Open("/etc/debian_version")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotIdentified
		}
		return nil, err
	}

	// use LSB to detect ubuntu, and linuxmint, debian does not have working lsb
	// by default.
	lsb, err := IdentifyLSBRelease()
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		return identifyDebianDerivertives(lsb)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &Info{
		Platform: PlatformDebian,
		Family:   FamilyDebian,
		Version:  string(bytes.TrimSpace(b)),
	}, nil
}

func identifyDebianDerivertives(lsb commands.Facts) (*Info, error) {
	info, ok := lsb.(*LSBInfo)
	if !ok {
		return nil, errors.New("failed to read lsb attributes")
	}
	platform := PlatformDebian
	switch {
	case strings.HasPrefix(info.ID, "Ubuntu"):
		platform = PlatformUbuntu
	case strings.HasPrefix(info.ID, "LinuxMint"):
		platform = PlatformLinuxMint
	}

	return &Info{
		Platform: platform,
		Family:   FamilyDebian,
		Version:  info.Release,
	}, nil
}
