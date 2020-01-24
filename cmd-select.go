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
	if args[0] == "_" {
		activeChatID = 0
		updateSidebar()
		outputView.SetOrigin(0, 0)
		outputView.SetCursor(0, 0)
		outputView.Clear()
		output(func(print printer) {
			print("Back to the machine room")
		})
		return nil
	}

	var user chat
	var count uint
	database.Model(&chat{}).Where("telegram_user_name = ?", args[0]).First(&user).Count(&count)
	if count == 0 {
		output(func(print printer) {
			print("No chat with that name present")
		})
	} else {
		activeContactID = 0
		activeChatID = user.ID
		if user.isLinked() {
			crew := user.FetchCrew()
			activeCrewID = crew.ID
		}
		outputView.SetOrigin(0, 0)
		outputView.SetCursor(0, 0)
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
