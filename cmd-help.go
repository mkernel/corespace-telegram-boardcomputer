package main

import (
	"fmt"
)

type helpCmd struct{}

func (helpCmd) Command() string {
	return "help"
}

func (helpCmd) Description() string {
	return "shows a list of available commands and executes their help"
}

func (helpCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "Usage: help [Command]")
	fmt.Fprintln(outputView, "without Command, lists all available commands")
}

func (helpCmd) Execute(args []string) error {
	if len(args) == 0 {
		//this is the default case.
		for _, cmd := range commands {
			fmt.Fprintf(outputView, "%s - %s\n", cmd.Command(), cmd.Description())
		}
	} else {
		command := args[0]
		for _, cmd := range commands {
			if cmd.Command() == command {
				cmd.Help(args[1:])
			}
		}
	}
	return nil
}
