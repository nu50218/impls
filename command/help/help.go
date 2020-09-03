package help

import (
	"flag"

	"github.com/nu50218/impls/command"
)

const name = "help"

var Command (command.Command) = &c{}

var flagSet = flag.NewFlagSet(name, flag.ExitOnError)

// オプション
var ()

type c struct{}

func (*c) Name() string {
	return name
}

func (*c) Run(args []string) error {
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	return nil
}
