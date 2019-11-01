package main

type botStatusCmd struct{}

func (botStatusCmd) Command() string {
	return "/status"
}

func (botStatusCmd) Description() string {
	return "Kurzzusammenfassung der Situation"
}

func (botStatusCmd) Execute(worker automationworker, args []string) {
	crew := worker.Chat.FetchCrew()
	worker.Chat.sendMessage(crew.Description)
}
