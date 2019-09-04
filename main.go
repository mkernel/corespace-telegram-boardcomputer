package main

import (
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jroimartin/gocui"
)

var outputView *gocui.View
var database *gorm.DB

func main() {
	filename := os.Args[1]
	database, _ = gorm.Open("sqlite3", filename)

	database.AutoMigrate(&globalSettings{})

	var settings globalSettings
	database.First(&settings)
	if database.NewRecord(settings) {
		//we have no dataset.
		database.Create(&settings)
	}

	defer database.Close()
	fillCommands()
	g, _ := gocui.NewGui(gocui.OutputNormal)
	g.Cursor = true
	g.SetManagerFunc(layout)
	g.SetKeybinding("commandline", gocui.KeyEnter, gocui.ModNone, operatecommand)

	g.MainLoop()
	g.Close()
}

func operatecommand(g *gocui.Gui, v *gocui.View) error {
	cmd := v.Buffer()
	cmd = strings.TrimSuffix(cmd, "\n")

	v.SetCursor(0, 0)
	v.Clear()
	params := strings.Split(cmd, " ")
	return run(params)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	outputView, _ = g.SetView("content", 0, 0, maxX-33, maxY-4)
	outputView.Wrap = true
	outputView.Autoscroll = true
	g.SetView("sidebar", maxX-32, 0, maxX-1, maxY-4)
	v, _ := g.SetView("commandline", 0, maxY-3, maxX-1, maxY-1)
	v.Editable = true
	g.SetCurrentView("commandline")
	return nil
}
