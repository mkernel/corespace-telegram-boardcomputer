package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var tgbot *tgbotapi.BotAPI

func setupTelegram() {
	var settings globalSettings
	database.First(&settings)
	if settings.APIKey != "" {
		var err error
		tgbot, err = tgbotapi.NewBotAPI(settings.APIKey)
		if err != nil {
			output <- "unable to connect to telegram"
		} else {
			output <- fmt.Sprintf("Telegram Identity: %s", tgbot.Self.UserName)
		}
		u := tgbotapi.NewUpdate(settings.APIOffset)
		u.Timeout = 60
		var botupdates tgbotapi.UpdatesChannel
		botupdates, err = tgbot.GetUpdatesChan(u)
		if err != nil {
			output <- "unable to establish updates connection"
		}
		go telegramFetcher(botupdates)
	} else {
		output <- "No Telegram API key set. SEt one and restart in order to start operations."
	}
}

func telegramFetcher(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		var settings globalSettings
		database.First(&settings)
		settings.APIOffset = update.UpdateID + 1
		database.Save(&settings)

		incomingmessage := update.Message
		userID := incomingmessage.From.ID

		storeduser := chat{TelegramID: userID}
		database.Where(&storeduser).First(&storeduser)
		if database.NewRecord(storeduser) {
			//this is a new record, so we have to fill in the details.
			storeduser.TelegramFirstName = incomingmessage.From.FirstName
			storeduser.TelegramLastName = incomingmessage.From.LastName
			storeduser.TelegramUserName = incomingmessage.From.UserName
			database.Create(&storeduser)
		}

		storedmessage := message{Inbound: true, Text: incomingmessage.Text, Date: incomingmessage.Date, ChatID: storeduser.ID, Read: false}
		database.Create(&storedmessage)

	}
}
