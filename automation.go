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
				worker, ok := automationqueues[item.Chat.ID]
				if !ok {
					worker = newWorker(item.Chat)
					automationqueues[item.Chat.ID] = worker
				}
				worker.Queue <- item.Message
			}
		}
	}
}

func handleUnlinkedMessage(item automationitem) {
	text := item.Message.Text
	var foundCrew crew
	database.Where(&crew{Code: text}).First(&foundCrew)
	if database.NewRecord(foundCrew) {
		item.Chat.sendMessage("Willkommen, Captain. Ich benötige Ihren Autorisierungscode, damit wir fortfahren können.")
	} else {
		foundCrew.Chat = item.Chat
		database.Save(&foundCrew)
		updateSidebar()
		item.Chat.sendMessage(fmt.Sprintf("Autorisierung bestätigt. Du bist der Captain der %s.", foundCrew.Name))
		item.Chat.sendMessage("Ich stehe dir jederzeit über eine Reihe von Befehlen zur Verfügung. Sende einfach '/help' für eine Liste.")
		item.Chat.sendMessage(foundCrew.Description)
		item.Chat.sendMessage(fmt.Sprintf("Du verfügst über %.2f AU", foundCrew.balance()))
		//TODO: check for unread messages from contacts and tell about those.
	}
}
