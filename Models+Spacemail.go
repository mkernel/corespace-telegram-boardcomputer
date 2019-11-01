package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (message *spacemail) print(view *gocui.View) {
	fmt.Fprintf(view, "%s\n", message.toString())
}

func (message *spacemail) toString() string {
	username := "<Crew>"
	if message.Inbound {
		username = "<Contact>"
	}
	result := fmt.Sprintf("%s: %s", username, message.Text)
	if message.Read == false && message.Inbound == false {
		message.Read = true
		database.Save(message)
	}
	return result
}
