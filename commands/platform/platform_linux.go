// Copyright 2015 MICHII Shunsuke. All rights reserved.

// +build linux

package platform

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/harukasan/orchestra-pit/commands"
	"github.com/harukasan/orchestra-pit/commands/exec"
)

// ErrNotIdentified is an error caused when the function could not identify the
// platform.
var ErrNotIdentified = errors.New("the platform could not identified")

// identifyFunc is a function which identifies the platform and retrieves the
// release information. If the platform could not identified, the function that
// implements identifyFunc should return ErrNotIdentified.
type identifyFunc func() (*Info, error)

var identifyFuncs = []identifyFunc{
	IdentifyDebianRelease,
	IdentifyRedHatRelease,
}

// Identify detects the platform and returns Info of the platform.
func Identify() (commands.Facts, error) {
	for _, f := range identifyFuncs {
		info, err := f()
		if err == nil {
			return info, nil
		}
		if err != ErrNotIdentified {
			return nil, err
		}
	}
	return nil, ErrNotIdentified
}

// IdentifyLSBRelease tires to retrieve the release information of LSB. If the
// platform is not supported for LSB, DetectLSBRelease returns nil.
func IdentifyLSBRelease() (commands.Facts, error) {
	lsbInfoCache.RLock()
	if lsbInfoCache.i != nil {
		defer lsbInfoCache.RUnlock()
		return lsbInfoCache.i, nil
	}
	lsbInfoCache.RUnlock()
	lsbInfoCache.Lock()
	defer lsbInfoCache.Unlock()

	info, err := readLSBFile()
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else if info != nil {
		lsbInfoCache.i = info
		return info, nil
	}

	info, err = execLSBRelease()
	if err != nil {
		return nil, err
	}
	lsbInfoCache.i = info
	return info, nil
}

var lsbInfoCache struct {
	sync.RWMutex
	i *LSBInfo
}

// readLSBFile reads the LSB release information from the lsb-release file.
func readLSBFile() (*LSBInfo, error) {
	file, err := os.Open("/etc/lsb-release")
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	parser := &LineParser{
		Delimiter:  '=',
		TrimSpaces: true,
		TrimQuotes: true,
	}
	m, err := parser.Parse(content)
	if err != nil {
		return nil, err
	}

	return &LSBInfo{
		ID:          string(m["DISTRIB_ID"]),
		Release:     string(m["DISTRIB_RELEASE"]),
		Codename:    string(m["CODENAME"]),
		Description: string(m["DESCRIPTION"]),
	}, nil
}

func execLSBRelease() (*LSBInfo, error) {
	out, err := exec.Command("/usr/bin/lsb_release -a").Output()
	if err != nil {
		return nil, err
	}

	parser := &LineParser{
		Delimiter:  ':',
		TrimSpaces: true,
	}
	m, err := parser.Parse(out)
	if err != nil {
		return nil, err
	}

	return &LSBInfo{
		ID:          string(m["Distributor ID"]),
		Release:     string(m["Release"]),
		Codename:    string(m["Codename"]),
		Description: string(m["Description"]),
	}, nil
}
