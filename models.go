package yuwiki

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
}

type Part struct {
	gorm.Model
}

type Page struct {
	gorm.Model
}
