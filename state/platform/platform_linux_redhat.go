// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package platform

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

// IdentifyRedHatRelease tires to identify the derivretives of Red Hat Linux.
// It supports following distributions:
//
//   - CentOS
//   - Red Hat Enterprise Linux
//
// If failed to identify the platform or the platform is the derivatives of
// Red Hat Linux, IdentifyRedHatRelease returns an ErrNotIdentifier.
func IdentifyRedHatRelease() (*Info, error) {
	file, err := os.Open("/etc/redhat-release")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotIdentified
		}
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var platform Name
	lb := bytes.ToLower(b)
	switch {
	case bytes.HasPrefix(lb, []byte("centos")):
		platform = PlatformCentOS
	case bytes.HasPrefix(lb, []byte("red hat enterprise")):
		platform = PlatformRHEL
	}

	version := ""
	if matches := regexp.MustCompile(`release (\d[\d.]*)`).FindSubmatch(b); len(matches) > 1 {
		version = string(matches[1])
	}

	return &Info{
		Platform: platform,
		Family:   FamilyRHEL,
		Version:  version,
	}, nil
}
