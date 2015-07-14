// Copyright 2015 MICHII Shunsuke. All rights reserved.
//

// +build darwin

package platform

import (
	"github.com/harukasan/orchestra-pit/state"
	"github.com/harukasan/orchestra-pit/state/exec"
)

// Identify detects the platform and returns Info of the platform.
func Identify() (state.Facts, error) {
	out, err := exec.Command("/usr/bin/sw_vers").Output()
	if err != nil {
		return nil, err
	}

	parser := &LineParser{
		Delimiter:  ':',
		TrimSpaces: true,
	}
	attrs, err := parser.Parse(out)
	if err != nil {
		return nil, err
	}

	return &Info{
		Platform:     PlatformOSX,
		Family:       FamilyOSX,
		Version:      string(attrs["ProductVersion"]),
		BuildVersion: string(attrs["BuildVersion"]),
	}, err
}
