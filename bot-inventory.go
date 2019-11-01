package main

import (
	"fmt"
	"strings"
)

type botInventoryCmd struct{}

func (botInventoryCmd) Command() string {
	return "/inventory"
}

func (botInventoryCmd) Description() string {
	return "Bericht zu Lagervorräten und Ausrüstung"
}

func (botInventoryCmd) Execute(worker automationworker, args []string) {
	items := worker.Chat.FetchCrew().fetchInventory()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Wir haben %d Gegenstände eingelagert.\n", len(items)))
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("* %s: %s\n", item.Name, item.Description))
	}
	builder.WriteString("Darüber hinaus sind unsere Nahrungsmittel- und Treibstoffvorräte ausreichend gefüllt.")
	worker.Chat.sendMessage(builder.String())
}
