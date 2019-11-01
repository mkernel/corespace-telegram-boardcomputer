package main

import (
	"fmt"
	"strings"
)

type botStatusCmd struct{}

func (botStatusCmd) Command() string {
	return "/status"
}

func (botStatusCmd) Description() string {
	return "Kurzzusammenfassung der Situation"
}

func (botStatusCmd) Execute(worker *automationworker, args []string) {
	crew := worker.Chat.FetchCrew()
	var builder strings.Builder
	builder.WriteString(crew.Description)
	builder.WriteString(fmt.Sprintf("\nDer Kontostand beträgt %.2f AU", crew.balance()))
	worker.Chat.sendMessage(builder.String())
}
