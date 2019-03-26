package yuwiki

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

var Db *gorm.DB

func DbConn() *gorm.DB {
	db, err := gorm.Open(Config.DataSource.Dialect, Config.DataSource.Url)
	if err != nil {
		log.Panic()
		return nil
	}
	return db
}

func setOwner(scope *gorm.Scope) {
	if scope.HasColumn("owner") {
		err := scope.SetColumn("owner", getUserId())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func InitDb(update bool) {
	dbFile := Config.DataSource.Url
	if Mkdirs(dbFile) {
		if _, err := os.OpenFile(Config.DataSource.Url, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			log.Fatal(err)
		}
	}
	Db = DbConn()
	Db.DB().SetMaxOpenConns(200)
	Db.DB().SetMaxIdleConns(100)
	Db.Callback().Create().Before("gorm:create").Register("set_owner", setOwner)
	Db.LogMode(false)
	book := &Book{}
	part := &Part{}
	page := &Page{}
	tag := &Tag{}
	sharedBook := &SharedBook{}
	if update {
		if Db.HasTable(book) {
			Db.AutoMigrate(book)
		} else {
			Db.CreateTable(book)
		}
		if Db.HasTable(part) {
			Db.AutoMigrate(part)
		} else {
			Db.CreateTable(part)
		}
		if Db.HasTable(page) {
			Db.AutoMigrate(page)
		} else {
			Db.CreateTable(page)
		}
		if Db.HasTable(tag) {
			Db.AutoMigrate(tag)
		} else {
			Db.CreateTable(tag)
		}
		if Db.HasTable(sharedBook) {
			Db.AutoMigrate(sharedBook)
		} else {
			Db.CreateTable(sharedBook)
		}
	} else {
		Db.DropTableIfExists(book, part, page, tag, sharedBook)
		Db.CreateTable(book, part, page, tag, sharedBook)
	}
}
