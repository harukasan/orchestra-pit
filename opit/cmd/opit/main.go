package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/harukasan/orchestra-pit/state/platform"
)

// Version specifies the release versioning string of opit command.
var Version = "(develop)"

type command interface {
	description() string
	run(arg []string) int
}

var commands = map[string]command{
	"apply":   applyCommand(),
	"test":    testCommand(),
	"version": &version{},
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	name := os.Args[1]
	if command := commands[name]; command != nil {
		os.Exit(command.run(os.Args[2:]))
	}
	fmt.Printf("unknown command: %s\n", name)
	usage()
	os.Exit(1)
}

func usage() {
	usage := `
Usage: opit COMMAND [ARGUMENTS]

opit is a command line interface of orchestra-pit.
`

	fmt.Println(strings.TrimSpace(usage))
	fmt.Println("\nCommands:")
	for name, command := range commands {
		fmt.Printf("  %-8s %s\n", name, command.description())
	}

	fmt.Println("")

	fmt.Println("Use \"opit COMMAND -h\" for more information about the commands.")

	p, err := platform.Identify()
	if err == nil {
		if platform.Family(p.Get("family")) == platform.FamilyDebian {
			fmt.Println("\n -- opit has Super Cow Powers now --")
		}
	}
}

func getCommandUsage(usage string, printDefaults func()) func() {
	return func() {
		fmt.Println(strings.TrimSpace(usage))
		fmt.Println("\nOptions:")
		printDefaults()
	}
}
