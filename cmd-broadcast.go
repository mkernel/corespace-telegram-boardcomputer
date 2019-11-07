package main

import (
	"fmt"
)

type broadcastCmd struct{}

func (broadcastCmd) Command() string {
	return "broadcast"
}

func (broadcastCmd) Description() string {
	return "sends a message to all chats that are linked to a crew."
}

func (broadcastCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (cmd broadcastCmd) Execute(args []string) error {
	var casted cmdlinesink = cmd
	inputfocus = &casted
	return nil
}

func (broadcastCmd) TextEntered(data string) error {
	var chats []chat
	database.Find(&chats);
	for _,chat := range chats {
		if(chat.isLinked()) {
			chat.sendMessage(data);
		}
	}
	return nil
}
