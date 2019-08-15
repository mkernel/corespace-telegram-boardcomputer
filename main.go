package main;

import (
	"github.com/jroimartin/gocui"
	"strings"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var OutputView *gocui.View;
var Database *gorm.DB;

func main() {
	filename := os.Args[1];
	Database, _ = gorm.Open("sqlite3",filename);

	Database.AutoMigrate(&GlobalSettings{});

	var settings GlobalSettings;
	Database.First(&settings);
	if(Database.NewRecord(settings)) {
		//we have no dataset.
		Database.Create(&settings);
	}

	defer Database.Close();
	Fill();
	g, _ := gocui.NewGui(gocui.OutputNormal);
	g.Cursor = true;
	g.SetManagerFunc(layout);
	g.SetKeybinding("commandline",gocui.KeyEnter, gocui.ModNone,operatecommand);

	g.MainLoop();
	g.Close();
}

func operatecommand(g *gocui.Gui, v *gocui.View) error {
	cmd := v.Buffer();
	cmd = strings.TrimSuffix(cmd,"\n");

	v.SetCursor(0,0);
	v.Clear();
	params := strings.Split(cmd," ");
	return Run(params);
}

func layout(g *gocui.Gui) error {
	maxX,maxY := g.Size();
	OutputView,_ = g.SetView("content",0,0,maxX-33,maxY-4);
	OutputView.Wrap = true;
	OutputView.Autoscroll = true;
	g.SetView("sidebar",maxX-32,0,maxX-1,maxY-4);
	v,_ := g.SetView("commandline",0,maxY-3,maxX-1,maxY-1);
	v.Editable = true;
	g.SetCurrentView("commandline");
	return nil;
}