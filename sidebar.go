package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func updateSidebar() {
	ui.Update(func(g *gocui.Gui) error {
		sidebar.Clear()
		fmt.Fprintln(sidebar, "Chats:")
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
		fmt.Fprintln(sidebar, "")
		fmt.Fprintln(sidebar, "Crews:")
		var crews []crew
		database.Find(&crews)
		for _, crew := range crews {
			current := " "
			if crew.isCurrent() {
				current = ">"
			}
			fmt.Fprintf(sidebar, "%s %s (%s)\n", current, crew.Name, crew.Code)
		}
		return nil
	})
}
