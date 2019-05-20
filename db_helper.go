package yuwiki

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var Db *gorm.DB

func DbConn() *gorm.DB {
	db, err := gorm.Open(Config.DataSource.Dialect, Config.DataSource.Url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

func setOwner(scope *gorm.Scope) {
	if scope.HasColumn("owner") {
		err := scope.SetColumn("owner", getUserId())
		if err != nil {
			log.Error(err)
		}
	}
}

func InitDb() {
	dbFile := Config.DataSource.Url
	if Mkdirs(dbFile) {
		if _, err := os.OpenFile(Config.DataSource.Url, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			log.Error(err)
		}
	}
	Db = DbConn()
	Db.DB().SetMaxOpenConns(200)
	Db.DB().SetMaxIdleConns(100)
	Db.Callback().Create().Before("gorm:create").Register("set_owner", setOwner)
	Db.LogMode(false)
	user := &User{}
	book := &Book{}
	part := &Part{}
	page := &Page{}
	historicalPage := &HistoricalPage{}
	tag := &Tag{}
	pageTag := &PageTag{}
	pageTemplate := &PageTemplate{}
	sharedBook := &SharedBook{}
	uploadFile := &UploadFile{}
	if Config.DataSource.DdlUpdate {
		if Db.HasTable(user) {
			Db.AutoMigrate(user)
		} else {
			Db.CreateTable(user)
		}
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
		if Db.HasTable(historicalPage) {
			Db.AutoMigrate(historicalPage)
		} else {
			Db.CreateTable(historicalPage)
		}
		if Db.HasTable(tag) {
			Db.AutoMigrate(tag)
		} else {
			Db.CreateTable(tag)
		}
		if Db.HasTable(pageTag) {
			Db.AutoMigrate(pageTag)
		} else {
			Db.CreateTable(pageTag)
		}
		if Db.HasTable(pageTemplate) {
			Db.AutoMigrate(pageTemplate)
		} else {
			Db.CreateTable(pageTemplate)
		}
		if Db.HasTable(sharedBook) {
			Db.AutoMigrate(sharedBook)
		} else {
			Db.CreateTable(sharedBook)
		}
		if Db.HasTable(uploadFile) {
			Db.AutoMigrate(uploadFile)
		} else {
			Db.CreateTable(uploadFile)
		}
	} else {
		Db.DropTableIfExists(user, book, part, page, historicalPage, tag, pageTag, pageTemplate, sharedBook, uploadFile)
		Db.CreateTable(user, book, part, page, historicalPage, tag, pageTag, pageTemplate, sharedBook, uploadFile)
	}
}
