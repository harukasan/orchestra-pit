package main

import (
	"flag"
	"os"
	"time"

	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/recipe"
	"github.com/harukasan/orchestra-pit/resource"
)

type test struct {
	// test command inherits apply command,
	// because test command have same options as apply command.
	*apply
}

func testCommand() *test {
	return &test{
		applyCommand(),
	}
}

func (c *test) description() string {
	return "test whether states of the host satisfies the given recipe file"
}

func (c *test) run(args []string) int {
	f := c.flags(args)
	c.initLogging()
	logger.Infof("Started at %s", time.Now().Format("2006-01-02T15:04:05-07:00"))

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	name := f.Arg(0)
	rec, err := recipe.ReadRecipe(name, wd)
	if err != nil {
		logger.Fatal(err)
	}

	exit := 0
	for _, res := range rec.Resources {
		logger.Debugf("------ testing %s", res)
		if err := resource.Test(res); err != nil {
			exit = 1
			logger.Errorf("[FAIL] %s", err)
			continue
		}
		logger.Infof("[ OK ] %s", res)
	}

	return exit
}

func (c *test) flags(args []string) *flag.FlagSet {
	usage := `
Usage: opit test [recipe]

Test whether states of the host satisfies the given recipe file.

The recipe must have the following extensions: yaml, yml, json, or rb.
If the recipe file is not specified, try to find the recipe file that is named
recipe.[ext] in current directory.
`

	f := flag.NewFlagSet("apply", flag.ExitOnError)
	f.Usage = getCommandUsage(usage, f.PrintDefaults)
	f.BoolVar(&c.DryRun, "dry-run", false, "report the commands that will have executed")
	c.loggingFlags(f)
	f.Parse(args)

	return f
}
