package main

import (
	"flag"
	"os"
	"time"

	"github.com/harukasan/orchestra-pit/opit/logger"
	"github.com/harukasan/orchestra-pit/opit/recipe"
	"github.com/harukasan/orchestra-pit/resource"
)

type apply struct {
	*logging
	DryRun bool
}

func applyCommand() *apply {
	return &apply{
		&logging{},
		false,
	}
}

func (c *apply) description() string {
	return "apply the recipe to the host"
}

func (c *apply) run(args []string) int {
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
			logger.Debugf("failed to test: %s", err)
		} else {
			logger.Infof("[ OK ] %s", res)
			continue
		}
		logger.Debugf("------ applying %s", res)
		if err := resource.Apply(res); err != nil {
			exit = 1
			logger.Debugf("failed to apply: %s", err)
			logger.Errorf("[FAIL] %s", res)
			continue
		}
		logger.Infof("[DONE] %s", res)
	}

	return exit
}

func (c *apply) flags(args []string) *flag.FlagSet {
	usage := `
Usage: opit apply [recipe]

Apply the recipe to the host.

The recipe must have the following extensions: yaml, yml, json, or rb.
If the recipe is not specified, try to find the recipe file named "recipe.[ext]"
in the current directory.
`

	f := flag.NewFlagSet("apply", flag.ExitOnError)
	f.Usage = getCommandUsage(usage, f.PrintDefaults)
	f.BoolVar(&c.DryRun, "dry-run", false, "report the commands that will have executed")
	c.loggingFlags(f)
	f.Parse(args)

	return f
}
