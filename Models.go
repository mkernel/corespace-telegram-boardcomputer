package main

import (
	"github.com/jinzhu/gorm"
)

type globalSettings struct {
	gorm.Model
	APIKey string
}
