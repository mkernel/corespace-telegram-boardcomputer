package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (message *spacemail) print(view *gocui.View) {
	fmt.Fprintf(view, "%s\n", message.toString())
}

func (message *spacemail) toString() string {
	username := "<??>"
	if message.Inbound {
		var c contact
		database.First(&c, message.ContactID)
		username = fmt.Sprintf("<%s>", c.Name)
	} else {
		var c crew
		database.First(&c, message.CrewID)
		username = fmt.Sprintf("<%s>", c.Name)
	}
	result := fmt.Sprintf("%s: %s", username, message.Text)
	if message.Read == false && message.Inbound == false {
		message.Read = true
		database.Save(message)
	}
	return result
}
