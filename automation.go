package main

import "fmt"

type automationitem struct {
	Chat    chat
	Message message
}

var automationqueue chan automationitem
var automationenabled map[uint]bool
var automationqueues map[uint]*automationworker

func setupAutomation() {
	automationenabled = make(map[uint]bool)
	automationqueues = make(map[uint]*automationworker)
	automationqueue = make(chan automationitem, 100)
	go automationWorker()
}

func automationWorker() {
	for item := range automationqueue {
		enabled, ok := automationenabled[item.Chat.ID]
		if !ok {
			automationenabled[item.Chat.ID] = true
			enabled = true
		}
		if enabled {
			if item.Chat.isLinked() == false {
				//this is an unlinked chat. we have to check things first.
				handleUnlinkedMessage(item)
			} else {
				//here we will handle linked accounts. As this got more complex we will spin off goroutines here, but not now.
				worker := getWorker(item.Chat.ID)
				worker.Queue <- item.Message
			}
		}
	}
}

func getWorker(chatid uint) *automationworker {
	worker, ok := automationqueues[chatid]
	if !ok {
		var chat chat
		database.First(&chat, chatid)
		worker = newWorker(chat)
		automationqueues[chatid] = worker
	}
	return worker
}

func handleUnlinkedMessage(item automationitem) {
	text := item.Message.Text
	var foundCrew crew
	database.Where(&crew{Code: text}).First(&foundCrew)

	var settings globalSettings
	database.First(&settings)

	if database.NewRecord(foundCrew) {
		item.Chat.sendMessage(settings.UnauthenticatedGreeting)
	} else {
		foundCrew.Chat = item.Chat
		database.Save(&foundCrew)
		updateSidebar()
		item.Chat.sendMessage(fmt.Sprintf(settings.AuthenticatedGreeting, foundCrew.Name))
		item.Chat.sendMessage(settings.AuthenticatedIntro)
		item.Chat.sendMessage(foundCrew.Description)
		item.Chat.sendMessage(fmt.Sprintf("Du verfügst über %.2f AU", foundCrew.balance()))
		var count uint
		database.Model(&spacemail{}).Where("crew_id = ? and inbound = ? and read = ?", foundCrew.ID, true, false).Count(&count)
		if count > 0 {
			item.Chat.sendMessage("Es gibt ungelesene Nachrichten.")
			//TODO: should we resend those?
		}
	}
}
