package main

func fillCommands() {
	commands = []command{
		exitCmd{},
		helpCmd{},
		apiKeyCmd{},
		selectCmd{},
		sendCmd{},
		crewCmd{},
		automationCmd{},
		unlinkCmd{},
		impersonateCmd{},
		membersCmd{},
		inventoryCmd{},
	}
}
