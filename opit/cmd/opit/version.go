package main

import (
	"fmt"
	"runtime"
)

type version struct{}

func (c *version) description() string {
	return "print version string"
}

func (c *version) run(arg []string) int {
	fmt.Printf("orchestra-pit version %s for %s\n", Version, runtime.GOOS)
	return 0
}
