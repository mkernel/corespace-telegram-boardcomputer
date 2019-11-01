package main

import (
	"fmt"
	"strings"
)

type botContactsCmd struct {
}

func (botContactsCmd) Command() string {
	return "/contacts"
}

func (botContactsCmd) Description() string {
	return "Zugriff auf die Kontaktliste"
}

func (botContactsCmd) Execute(worker *automationworker, args []string) {
	contacts := worker.Chat.FetchCrew().fetchContacts()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Es sind %d Kontakte in der Datenbank:\n", len(contacts)))
	for _, contact := range contacts {
		builder.WriteString(fmt.Sprintf("* %s (%d ungelesene)\n", contact.Name, contact.numCrewUnread()))
	}
	builder.WriteString("Mit /info NAME kannst du mehr über jeden Kontakt erfahren.\n")
	builder.WriteString("Mit /read NAME kannst du auf die eingegangenen Nachrichten zugreifen.\n")
	builder.WriteString("Mit /write NAME kannst du eine Nachricht schreiben.")
	builder.WriteString("Aber es geht natürlich auch weniger förmlich.")
	worker.Chat.sendMessage(builder.String())
}
