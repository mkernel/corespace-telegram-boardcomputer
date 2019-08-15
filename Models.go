package main;

import (
	"github.com/jinzhu/gorm"
)

type GlobalSettings struct {
	gorm.Model
	APIKey string;
}