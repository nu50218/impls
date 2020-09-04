package help

import (
	"flag"
	"fmt"

	"github.com/nu50218/impls/command"
)

const name = "help"

var flagSet = flag.NewFlagSet(name, flag.ExitOnError)

// オプション
var ()

type c struct {
	commands []command.Command
}

func New(commands ...command.Command) command.Command {
	return &c{
		commands: commands,
	}
}

func (*c) Name() string {
	return name
}

func (*c) Description() string {
	return "help"
}

func (*c) FlagSet() *flag.FlagSet {
	return flagSet
}

func (cc *c) Run(args []string) error {
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// $ go help [subcommand]
	if len(args) != 0 {
		for _, command := range cc.commands {
			if args[0] == command.Name() {
				command.FlagSet().Usage()
				return nil
			}
		}
		return fmt.Errorf("%s not found", args[0])
	}

	fmt.Println("ʕ◔ϖ◔ʔ　impls")
	fmt.Println()

	fmt.Println("ʕ◔ϖ◔ʔ　see help for each command by $ impls help [subcommand]")
	fmt.Println()

	fmt.Println("subcommands:")
	for i, command := range cc.commands {
		if i != 0 {
			fmt.Println()
		}
		fmt.Printf("  - %s\n      %s\n", command.Name(), command.Description())
	}

	return nil
}
