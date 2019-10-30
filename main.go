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
var output chan string
var ui *gocui.Gui
var inputfocus *cmdlinesink

var activeChatID uint

func main() {
	output = make(chan string, 100)
	filename := os.Args[1]
	database, _ = gorm.Open("sqlite3", filename)

	setupDatabase(database)
	setupTelegram()

	defer database.Close()
	fillCommands()
	g, _ := gocui.NewGui(gocui.OutputNormal)
	ui = g
	g.Cursor = true
	g.SetManagerFunc(layout)
	g.SetKeybinding("commandline", gocui.KeyEnter, gocui.ModNone, operatecommand)
	go writeOutput(g)
	g.MainLoop()
	g.Close()
}

func writeOutput(g *gocui.Gui) {
	for msg := range output {
		g.Update(func(g *gocui.Gui) error {
			fmt.Fprintln(outputView, msg)
			return nil
		})
	}
}

func operatecommand(g *gocui.Gui, v *gocui.View) error {
	cmd := v.Buffer()
	cmd = strings.TrimSuffix(cmd, "\n")

	v.SetCursor(0, 0)
	v.Clear()
	if inputfocus != nil {
		result := (*inputfocus).TextEntered(cmd)
		inputfocus = nil
		return result
	}
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
