package main

import (
	"fmt"
	"time"
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
		var user chat
		database.First(&user, activeChatID)
		user.sendMessage(data)
	} else if activeContactID != 0 {
		//TODO: implement linked crews
		dataset := spacemail{CrewID: activeCrewID, ContactID: activeContactID, Inbound: true, Date: int(time.Now().Unix()), Text: data, Read: false}
		database.Create(&dataset)
		output <- dataset.toString()
		var crew crew
		database.Preload("Chat").First(&crew, activeCrewID)
		if crew.ChatID != 0 {
			var contact contact
			database.First(&contact, activeContactID)
			crew.Chat.sendMessage(fmt.Sprintf("Ich habe eine Nachricht von %s empfangen.", contact.Name))
		}
	}
	return nil
}
