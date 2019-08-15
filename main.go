package main;

import (
	"github.com/jroimartin/gocui"
	"strings"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var OutputView *gocui.View;

func main() {
	filename := os.Args[1];
	db, _ := gorm.Open("sqlite3",filename);

	db.AutoMigrate(&GlobalSettings{});

	var settings GlobalSettings;
	db.First(&settings);
	if(db.NewRecord(settings)) {
		//we have no dataset.
		db.Create(&settings);
	}

	defer db.Close();
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
	g.SetView("sidebar",maxX-32,0,maxX-1,maxY-4);
	v,_ := g.SetView("commandline",0,maxY-3,maxX-1,maxY-1);
	v.Editable = true;
	g.SetCurrentView("commandline");
	return nil;
}