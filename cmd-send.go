package main

import (
	"fmt"
)

type sendCmd struct{}

func (sendCmd) Command() string {
	return "send"
}

func (sendCmd) Description() string {
	return "Sends a board computer message"
}

func (sendCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (sendcommand sendCmd) Execute(args []string) error {
	if activeChatID == 0 {
		output <- "You can only send a message with an active chat"
	} else {
		var casted cmdlinesink = sendcommand
		inputfocus = &casted
	}
	return nil
}

func (sendCmd) TextEntered(data string) error {
	var user chat
	database.First(&user, activeChatID)
	user.sendMessage(data)
	return nil
}
