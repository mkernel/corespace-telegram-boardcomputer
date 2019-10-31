package main

import (
	"fmt"
)

type unlinkCmd struct{}

func (unlinkCmd) Command() string {
	return "unlink"
}

func (unlinkCmd) Description() string {
	return "unlinks the active chat from its crew"
}

func (unlinkCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (unlinkCmd) Execute(args []string) error {
	var chat chat
	database.First(&chat, activeChatID)
	crew := chat.FetchCrew()
	crew.ChatID = 0
	database.Save(&crew)
	updateSidebar()
	return nil
}
