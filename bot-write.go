package main

import (
	"strings"
)

type botWriteCmd struct {
	ContactID uint
}

func (botWriteCmd) Command() string {
	return "/write"
}

func (botWriteCmd) Description() string {
	return "Nachricht an einen Kontakt senden (Aufruf: /write NAME)"
}

func (cmd botWriteCmd) Execute(worker *automationworker, args []string) {
	if len(args) == 0 {
		worker.Chat.sendMessage("Keinen Kontaktnamen angegeben")
		return
	}
	filter := contact{OwnerID: worker.Chat.FetchCrew().ID}
	var contacts []contact
	database.Where(&filter).Find(&contacts)
	var contact contact
	found := false
	for _, item := range contacts {
		if strings.ToLower(item.Name) == strings.ToLower(args[0]) {
			contact = item
			found = true
		}
	}
	if found == false {
		worker.Chat.sendMessage("Keinen passenden Kontakt in der Datenbank gefunden.")
		return
	}
	cmd.ContactID = contact.ID
	if len(args) > 1 {
		message := args[1:]
		text := strings.Join(message, " ")
		contact.sendMessageToContact(text)
		worker.Chat.sendMessage("Ich habe die Nachricht übertragen.")
	} else {
		worker.Chat.sendMessage("Ich bin ganz Ohr.")
		var casted botDataSink = cmd
		worker.CurrentFocus = &casted
	}
}

func (cmd botWriteCmd) OnMessage(worker *automationworker, msg message) {
	var destcontact contact
	database.First(&destcontact, cmd.ContactID)
	destcontact.sendMessageToContact(msg.Text)
	updateSidebar()
	worker.Chat.sendMessage("Ich habe die Nachricht übertragen.")
}
