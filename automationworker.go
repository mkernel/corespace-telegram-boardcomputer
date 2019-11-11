package main

import "strings"

type automationworker struct {
	Chat         chat
	Queue        chan message
	Commands     []botCommand
	CurrentFocus *botDataSink
}

func newWorker(chat chat) *automationworker {
	worker := &automationworker{Chat: chat}
	worker.start()
	return worker
}

func (worker *automationworker) start() {
	worker.Queue = make(chan message, 100)
	worker.setDefaultCommandSet()
	go worker.work()
}

func (worker *automationworker) setDefaultCommandSet() {
	worker.Commands = []botCommand{
		botHelpCmd{},
		botCrewCmd{},
		botStatusCmd{},
		botInventoryCmd{},
		botContactsCmd{},
		botInfoCmd{},
		botWriteCmd{},
		botConnectCmd{},
		botCallCmd{},
	}
}

func (worker *automationworker) setCommandSet(commands []botCommand) {
	worker.Commands = commands
}

func (worker *automationworker) work() {
	for message := range worker.Queue {
		if message.Text[0] == '/' {
			//it starts with / so it will be a command.
			message.Read = true
			database.Save(&message)
			updateSidebar()
			args := strings.Split(message.Text, " ")
			found := false
			for _, cmd := range worker.Commands {
				if cmd.Command() == args[0] {
					cmd.Execute(worker, args[1:])
					found = true
					break
				}
			}
			if !found {
				worker.Chat.sendMessage("Den Befehl habe ich leider nicht verstanden. Versuche es einmal mit /help")
			}
		} else if worker.CurrentFocus != nil {
			backup := worker.CurrentFocus
			worker.CurrentFocus = nil
			(*backup).OnMessage(worker, message)
			message.Read = true
			database.Save(&message)
			updateSidebar()
		} else if conferences.isCrewInOngoingCall(worker.Chat.FetchCrew()) {
			conferences.transmitFromCrew(worker.Chat.FetchCrew(), message.Text)
		}
	}
}
