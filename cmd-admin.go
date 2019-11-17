package main

import (
	"fmt"
)

type adminCmd struct{}

func (adminCmd) Command() string {
	return "admin"
}

func (adminCmd) Description() string {
	return "disables or enables admin mode for the active chat"
}

func (adminCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (adminCmd) Execute(args []string) error {
	if len(args) != 1 || activeChatID == 0 {
		return nil
	}
	var chat chat
	database.First(&chat, activeChatID)
	if args[0] == "enable" {
		chat.Admin = true
		output <- "(admin enabled)"
	} else if args[0] == "disable" {
		chat.Admin = false
		output <- "(admin disabled)"
	}
	database.Save(&chat)
	return nil
}
