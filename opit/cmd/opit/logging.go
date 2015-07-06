package main

import (
	"flag"
	"os"

	"github.com/harukasan/orchestra-pit/opit/logger"
)

type logging struct {
	Verbose     bool
	Quiet       bool
	TextLogFile string
	JSONLogFile string
}

func (c *logging) loggingFlags(f *flag.FlagSet) {
	f.StringVar(&c.TextLogFile, "log", "", "output logs to the file in the same format as stdout")
	f.StringVar(&c.JSONLogFile, "log-json", "", "output logs to the file in JSON format")
	f.BoolVar(&c.Quiet, "q", false, "omit to output information to stdout, run silently")
	f.BoolVar(&c.Verbose, "v", false, "print messages verbosity")
}

func (c *logging) initLogging() {
	if c.Quiet {
		logger.RemoveOutput(logger.StdoutOutput)
	}
	if c.TextLogFile != "" {
		file, err := os.OpenFile(c.TextLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Fatalf("can not open the log file, %s", err)
		}
		logger.AddOutput(logger.TextOutput(file))
	}
	if c.JSONLogFile != "" {
		file, err := os.OpenFile(c.JSONLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Fatalf("can not open the JSON log file, %s", err)
		}
		logger.AddOutput(logger.JSONOutput(file))
	}
}
