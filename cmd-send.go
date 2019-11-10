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
	if activeChatID == 0 && activeContactID == 0 {
		output <- "You can only send a message with an active chat or contact"
	} else {
		var casted cmdlinesink = sendcommand
		inputfocus = &casted
	}
	return nil
}

func (sendCmd) TextEntered(data string) error {
	if activeChatID != 0 {
		//TODO: check if on call and send to conference system
		var user chat
		database.First(&user, activeChatID)
		user.sendMessage(data)
	} else if activeContactID != 0 {
		//technical, we should opt for "not for linked contacts"
		var destination contact
		database.First(&destination, activeContactID)
		destination.sendMessageToCrew(data)
	}
	return nil
}
