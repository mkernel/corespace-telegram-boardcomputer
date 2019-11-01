package main

import "time"

type botWriteCmd struct {
	ContactID uint
}

func (botWriteCmd) Command() string {
	return "/write"
}

func (botWriteCmd) Description() string {
	return "Nachricht an einen Kontakt senden"
}

func (cmd botWriteCmd) Execute(worker *automationworker, args []string) {
	filter := contact{CrewID: worker.Chat.FetchCrew().ID, Name: args[0]}
	var found contact
	database.Where(&filter).First(&found)
	cmd.ContactID = found.ID
	worker.Chat.sendMessage("Ich bin ganz Ohr.")
	var casted botDataSink = cmd
	worker.CurrentFocus = &casted
}

func (cmd botWriteCmd) OnMessage(worker automationworker, msg message) {
	//TODO: support linked contacts
	var contact contact
	database.First(&contact, cmd.ContactID)
	spacemail := spacemail{CrewID: contact.CrewID, ContactID: cmd.ContactID, Inbound: false, Read: false, Date: int(time.Now().Unix()), Text: msg.Text}
	database.Create(&spacemail)
	if activeContactID == cmd.ContactID {
		output <- spacemail.toString()
	}
	updateSidebar()
	worker.Chat.sendMessage("Ich habe die Nachricht übertragen.")
}
