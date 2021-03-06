package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nu50218/impls/command"
	"github.com/nu50218/impls/command/help"
	"github.com/nu50218/impls/command/iface"
	"github.com/nu50218/impls/command/types"
	"github.com/nu50218/impls/command/vars"
)

func run(args []string) error {
	if len(args) < 1 {
		return errors.New("you must pass a subcommand")
	}

	subCommand := args[0]
	subCommands := []command.Command{
		iface.Command,
		types.Command,
		vars.Command,
	}

	subCommands = append(subCommands, help.New(subCommands...))

	for _, sc := range subCommands {
		if sc.Name() == subCommand {
			return sc.Run(args[1:])
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subCommand)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
