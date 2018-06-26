package main

import (
	"fmt"
	"io"
	"os"

	"github.com/satococoa/prbot/lib"
)

const (
	// ExitCodeOK success
	ExitCodeOK = iota
	// ExitCodeSettingError setting error
	ExitCodeSettingError
	// ExitCodeExecutionError execution error
	ExitCodeExecutionError
)

// CLI command line interface
type CLI struct{ outStream, errStream io.Writer }

// Run run
func (c *CLI) Run(args []string) int {
	setting, err := prbot.NewSetting()
	if err != nil {
		fmt.Print(c.errStream, err)
		return ExitCodeSettingError
	}

	err = prbot.Execute(setting)
	if err != nil {
		fmt.Print(c.errStream, err)
		return ExitCodeExecutionError
	}

	return ExitCodeOK
}

func main() {
	cli := &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	os.Exit(cli.Run(os.Args))
}
