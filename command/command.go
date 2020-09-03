package command

type Command interface {
	Name() string
	Run([]string) error
}
