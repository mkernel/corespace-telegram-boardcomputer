package main

type botInfoCmd struct{}

func (botInfoCmd) Command() string {
	return "/info"
}

func (botInfoCmd) Description() string {
	return "Informationen zu einem Kontakt abrufen"
}

func (botInfoCmd) Execute(worker *automationworker, args []string) {
	filter := contact{CrewID: worker.Chat.FetchCrew().ID, Name: args[0]}
	var found contact
	database.Where(&filter).First(&found)
	worker.Chat.sendMessage(found.Description)
}
