package main

import (
	"fmt"
)

type impersonateCmd struct{}

func (impersonateCmd) Command() string {
	return "impersonate"
}

func (impersonateCmd) Description() string {
	return "executes a message as if it came from the user"
}

func (impersonateCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (command impersonateCmd) Execute(args []string) error {
	var cast cmdlinesink = command
	inputfocus = &cast
	return nil
}

func (command impersonateCmd) TextEntered(data string) error {
	var user chat
	database.First(&user, activeChatID)
	simulatedmessage := message{Text: data}
	automationqueue <- automationitem{Chat: user, Message: simulatedmessage}
	return nil
}
