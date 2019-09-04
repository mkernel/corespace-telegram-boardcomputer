package main

func fillCommands() {
	commands = []command{
		exitCmd{},
		helpCmd{},
		apiKeyCmd{},
	}
}
