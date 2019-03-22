package yuwiki

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"strings"
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
	if scope.HasColumn("Owner") {
		err := scope.SetColumn("Owner", GetUserId)
		if err != nil {
			log.Fatal("设置所有者失败", err)
		}
	}
}

func InitDb(update bool) {
	dbFile := Config.DataSource.Url
	if _, err := os.Open(dbFile); err != nil {
		dir := string(dbFile[:strings.LastIndex(dbFile, "/")])
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal(fmt.Sprintf("创建数据库文件所在目录[%s]失败", dir), err)
		} else if _, err := os.Create(Config.DataSource.Url); err != nil {
			log.Fatal(fmt.Sprintf("创建数据库文件[%s]失败", dbFile), err)
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
