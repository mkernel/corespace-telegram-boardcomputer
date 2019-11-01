package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type botCrewCmd struct{}

func (botCrewCmd) Command() string {
	return "/crew"
}

func (botCrewCmd) Description() string {
	return "Zugriff auf die Akten der Mannschaft"
}

func (botCrewCmd) Execute(worker automationworker, args []string) {
	var members []member
	crew := worker.Chat.FetchCrew()
	filter := member{CrewID: crew.ID}
	database.Where(&filter).Find(&members)
	worker.Chat.sendMessage(fmt.Sprintf("Deine Crew besteht aus %d Mitgliedern.", len(members)))
	for _, member := range members {
		filename := member.Filename()
		msg := tgbotapi.NewPhotoUpload(worker.Chat.TelegramChatID, filename)
		tgbot.Send(msg)
	}
	worker.Chat.sendMessage("Fertig.")
}
