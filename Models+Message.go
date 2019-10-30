package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (message *message) print(view *gocui.View) {
	fmt.Fprintf(view, "%s\n", message.toString())
}

func (message *message) toString() string {
	username := "<Bot>"
	if message.Inbound {
		username = "<User>"
	}
	result := fmt.Sprintf("%s: %s", username, message.Text)
	if message.Read == false && message.Inbound {
		message.Read = true
		database.Save(message)
	}
	return result
}
