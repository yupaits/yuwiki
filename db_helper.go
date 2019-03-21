package yuwiki

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"strings"
)

func DbConn() *gorm.DB {
	db, err := gorm.Open(Config.DataSource.Dialect, Config.DataSource.Url)
	if err != nil {
		log.Panic()
		return nil
	}
	db.Callback().Create().Before("gorm:create").Register("set_owner", setOwner)
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
	db := DbConn()
	defer db.Close()
	tables := []interface{}{&Book{}, &Part{}, &Page{}, &Tag{}}
	if update {
		var updateTables []interface{}
		var createTables []interface{}
		for _, table := range tables {
			if db.HasTable(table) {
				updateTables = append(updateTables, table)
			} else {
				createTables = append(createTables, table)
			}
		}
		db.AutoMigrate(updateTables)
		db.CreateTable(createTables)
	} else {
		db.DropTableIfExists(tables)
		db.CreateTable(tables)
	}
}
