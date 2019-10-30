package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func updateSidebar() {
	ui.Update(func(g *gocui.Gui) error {
		sidebar.Clear()
		var users []chat
		database.Order("telegram_user_name asc").Find(&users)
		for _, user := range users {
			marker := ""
			if user.hasUnread() {
				marker = "*"
			}
			linked := "?"
			if user.isLinked() {
				linked = ""
			}
			current := " "
			if user.isCurrent() {
				current = ">"
			}
			fmt.Fprintf(sidebar, "%s %s%s%s\n", current, user.TelegramUserName, linked, marker)
		}
		return nil
	})
}
