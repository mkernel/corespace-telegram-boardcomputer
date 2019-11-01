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
	filter := contact{OwnerID: worker.Chat.FetchCrew().ID, Name: args[0]}
	var found contact
	database.Where(&filter).First(&found)
	cmd.ContactID = found.ID
	worker.Chat.sendMessage("Ich bin ganz Ohr.")
	var casted botDataSink = cmd
	worker.CurrentFocus = &casted
}

func (cmd botWriteCmd) OnMessage(worker automationworker, msg message) {
	var destcontact contact
	database.First(&destcontact, cmd.ContactID)
	mail := spacemail{CrewID: destcontact.OwnerID, ContactID: cmd.ContactID, Inbound: false, Read: false, Date: int(time.Now().Unix()), Text: msg.Text}
	database.Create(&mail)
	if destcontact.CrewID != 0 {
		//this is a linked contact. So we have to duplicate the spacemail over to the other crew.
		var mirrorcontact contact
		database.Where(&contact{OwnerID: destcontact.CrewID, CrewID: destcontact.OwnerID}).First(&mirrorcontact)
		mirrormail := spacemail{CrewID: mirrorcontact.OwnerID, ContactID: mirrorcontact.ID, Inbound: true, Read: false, Date: int(time.Now().Unix()), Text: msg.Text}
		database.Create(&mirrormail)
		if activeContactID == mirrorcontact.ID {
			output <- mail.toString()
		}
	}
	if activeContactID == cmd.ContactID {
		output <- mail.toString()
	}
	updateSidebar()
	worker.Chat.sendMessage("Ich habe die Nachricht übertragen.")
}
