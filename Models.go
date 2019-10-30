package main

import (
	"github.com/jinzhu/gorm"
)

type globalSettings struct {
	gorm.Model
	APIKey    string
	APIOffset int
}

type chat struct {
	gorm.Model
	TelegramID        int
	TelegramFirstName string
	TelegramLastName  string
	TelegramUserName  string
	Messages          []message
}

type message struct {
	gorm.Model
	Inbound bool
	Text    string `gorm:"type:text"`
	Date    int
	ChatID  uint
	Read    bool
}

type crew struct {
	gorm.Model
	Code         string
	Name         string
	Description  string `gorm:"type:text"`
	ChatID       uint
	Chat         chat
	Members      []member
	Items        []item
	Transactions []transaction
	Contacts     []contact `gorm:"foreignkey:OwnerID"`
	Spacemails   []spacemail
}

type member struct {
	gorm.Model
	Name   string
	CrewID uint
}

type item struct {
	gorm.Model
	Name        string
	Description string `gorm:"type:text"`
	CrewID      uint
}

type transaction struct {
	gorm.Model
	Date        int
	Value       float64
	CrewID      uint
	Description string
}

type contact struct {
	gorm.Model
	OwnerID     uint
	CrewID      uint
	Crew        crew
	Name        string
	Description string `gorm:"type:text"`
	Spacemails  []spacemail
}

type spacemail struct {
	gorm.Model
	Text      string `gorm:"type:text"`
	Date      int
	Read      bool
	CrewID    uint
	Inbound   bool
	ContactID uint
}

func setupDatabase(db *gorm.DB) {
	database.AutoMigrate(&globalSettings{})
	database.AutoMigrate(&chat{})
	database.AutoMigrate(&message{})
	database.AutoMigrate(&crew{})
	database.AutoMigrate(&member{})
	database.AutoMigrate(&item{})
	database.AutoMigrate(&transaction{})
	database.AutoMigrate(&contact{})
	database.AutoMigrate(&spacemail{})

	var settings globalSettings
	database.First(&settings)
	if database.NewRecord(settings) {
		//we have no dataset.
		database.Create(&settings)
	}
}
