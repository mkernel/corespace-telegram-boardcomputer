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
	Text    string
	Date    int
	ChatID  uint
}
