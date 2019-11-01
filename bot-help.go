package main

import (
	"fmt"
	"strings"
)

type botHelpCmd struct{}

func (botHelpCmd) Command() string {
	return "/help"
}

func (botHelpCmd) Description() string {
	return "zeigt diese Liste an."
}

func (botHelpCmd) Execute(worker automationworker, args []string) {
	var builder strings.Builder
	for _, cmd := range worker.Commands {
		builder.WriteString(fmt.Sprintf("%s - %s\n", cmd.Command(), cmd.Description()))
	}
	worker.Chat.sendMessage(builder.String())
}
