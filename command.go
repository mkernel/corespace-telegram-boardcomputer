package main

type command interface {
	Command() string
	Description() string
	Help(args []string)
	Execute(args []string) error
}
