package main

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
	filter := contact{OwnerID: worker.Chat.FetchCrew().ID, Name: args[0]}
	var found contact
	database.Where(&filter).First(&found)
	if database.NewRecord(&found) {
		worker.Chat.sendMessage("Keinen passenden Kontakt in der Datenbank gefunden.")
		return
	}
	cmd.ContactID = found.ID
	worker.Chat.sendMessage("Ich bin ganz Ohr.")
	var casted botDataSink = cmd
	worker.CurrentFocus = &casted
}

func (cmd botWriteCmd) OnMessage(worker automationworker, msg message) {
	var destcontact contact
	database.First(&destcontact, cmd.ContactID)
	destcontact.sendMessageToContact(msg.Text)
	updateSidebar()
	worker.Chat.sendMessage("Ich habe die Nachricht übertragen.")
}
