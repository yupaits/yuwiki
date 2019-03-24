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
		err := scope.SetColumn("owner", GetUserId())
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
	if update {
		if Db.HasTable(&Book{}) {
			Db.AutoMigrate(&Book{})
		} else {
			Db.CreateTable(&Book{})
		}
		if Db.HasTable(&Part{}) {
			Db.AutoMigrate(&Part{})
		} else {
			Db.CreateTable(&Part{})
		}
		if Db.HasTable(&Page{}) {
			Db.AutoMigrate(&Page{})
		} else {
			Db.CreateTable(&Page{})
		}
		if Db.HasTable(&Tag{}) {
			Db.AutoMigrate(&Tag{})
		} else {
			Db.CreateTable(&Tag{})
		}
	} else {
		Db.DropTableIfExists(&Book{}, &Part{}, &Page{}, &Tag{})
		Db.CreateTable(&Book{}, &Part{}, &Page{}, &Tag{})
	}
}
