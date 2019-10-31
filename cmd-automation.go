package main

import (
	"fmt"
)

type automationCmd struct{}

func (automationCmd) Command() string {
	return "automation"
}

func (automationCmd) Description() string {
	return "disables or enables automation for the active chat"
}

func (automationCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (automationCmd) Execute(args []string) error {
	if len(args) != 1 {
		return nil
	}
	if args[0] == "enable" {
		automationenabled[activeChatID] = true
		output <- "(automation enabled)"
	} else if args[0] == "disable" {
		automationenabled[activeChatID] = false
		output <- "(automation disabled)"
	}
	return nil
}
