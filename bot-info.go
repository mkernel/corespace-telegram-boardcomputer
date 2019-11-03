package main

type botInfoCmd struct{}

func (botInfoCmd) Command() string {
	return "/info"
}

func (botInfoCmd) Description() string {
	return "Informationen zu einem Kontakt abrufen (Aufruf: /info NAME)"
}

func (botInfoCmd) Execute(worker *automationworker, args []string) {
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
	worker.Chat.sendMessage(found.Description)
}
