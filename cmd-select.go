package main

import (
	"fmt"
)

type selectCmd struct{}

func (selectCmd) Command() string {
	return "select"
}

func (selectCmd) Description() string {
	return "switches the current chat"
}

func (selectCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (selectCmd) Execute(args []string) error {
	if len(args) != 1 {
		return nil
	}
	var user chat
	var count uint
	database.Model(&chat{}).Where("telegram_user_name = ?", args[0]).First(&user).Count(&count)
	if count == 0 {
		output <- "No chat with that name present"
	} else {
		activeChatID = user.ID
		outputView.Clear() //when switching, we clear the output.
		messages := user.fetchMessages()
		lastOne := len(messages)
		firstOne := lastOne - 50
		if firstOne < 0 {
			firstOne = 0
		}
		toPrint := messages[firstOne:lastOne]
		for _, msg := range toPrint {
			msg.print(outputView)
		}
		updateSidebar()
	}
	return nil
}
