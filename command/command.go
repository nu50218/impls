package command

import "flag"

type Command interface {
	Name() string
	Description() string
	FlagSet() *flag.FlagSet
	Run([]string) error
}
