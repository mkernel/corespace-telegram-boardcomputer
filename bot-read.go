package main

import "fmt"

type botReadCmd struct{}

func (botReadCmd) Command() string {
	return "/read"
}

func (botReadCmd) Description() string {
	return "Zugriff auf die eingegangenen Nachrichten eines Kontakts (Aufruf: /read NAME)"
}

func (botReadCmd) Execute(worker *automationworker, args []string) {
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
	var messages []spacemail
	database.Where("crew_id = ? and contact_id = ? and inbound = ? and read = ?", worker.Chat.FetchCrew().ID, found.ID, true, false).Order("date asc").Find(&messages)
	if len(messages) == 0 {
		worker.Chat.sendMessage("Es gibt keine neuen Nachrichten.")
	}
	for _, message := range messages {
		worker.Chat.sendMessage(fmt.Sprintf("<%s> %s", args[0], message.Text))
		message.Read = true
		database.Save(&message)
	}
}
