package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jroimartin/gocui"
)

var outputView *gocui.View
var sidebar *gocui.View
var database *gorm.DB
var ui *gocui.Gui
var inputfocus *cmdlinesink

var activeChatID uint
var activeCrewID uint
var activeContactID uint

func main() {
	filename := os.Args[1]
	database, _ = gorm.Open("sqlite3", filename)
	g, _ := gocui.NewGui(gocui.OutputNormal)
	ui = g
	g.Cursor = true
	g.SetManagerFunc(layout)
	g.SetKeybinding("commandline", gocui.KeyEnter, gocui.ModNone, operatecommand)
	setupDatabase(database)
	setupAutomation()
	setupTelegram()
	SetupAdmin()
	setupHTTP()

	defer database.Close()
	fillCommands()
	//before we start, we have to disconnect everyone from after being connected during a crash.
	var chats []chat
	database.Where(&chat{OpenConnection: true}).Find(&chats)
	for _, chat := range chats {
		chat.OpenConnection = false
		database.Save(&chat)
		chat.sendMessage("Die Verbindung wurde unterbrochen")
	}
	g.MainLoop()
	g.Close()
}

type printer func(string)

func output(context func(printer printer)) {
	ui.Update(func(g *gocui.Gui) error {
		context(func(line string) {
			fmt.Fprintln(outputView, line)
		})
		return nil
	})
}

func operatecommand(g *gocui.Gui, v *gocui.View) error {
	cmd := v.Buffer()
	cmd = strings.TrimSuffix(cmd, "\n")
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	v.Clear()
	if inputfocus != nil {
		instance := inputfocus
		inputfocus = nil
		return (*instance).TextEntered(cmd)
	}
	fmt.Fprintf(outputView, "> %s\n", cmd)
	params := strings.Split(cmd, " ")
	return run(params)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	outputView, _ = g.SetView("content", 0, 0, maxX-33, maxY-4)
	outputView.Wrap = true
	outputView.Autoscroll = true
	initialupdate := sidebar == nil
	sidebar, _ = g.SetView("sidebar", maxX-32, 0, maxX-1, maxY-4)
	v, _ := g.SetView("commandline", 0, maxY-3, maxX-1, maxY-1)
	v.Editable = true
	g.SetCurrentView("commandline")
	if initialupdate {
		updateSidebar()
	}
	return nil
}
