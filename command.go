package main


type Command interface {
	Command() string
	Description() string
	Help(args []string)
	Execute(args []string) error
}
