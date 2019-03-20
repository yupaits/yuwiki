package yuwiki

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func DbConn() *gorm.DB {
	db, err := gorm.Open(Config.DataSource.Dialect, Config.DataSource.Url)
	if err != nil {
		log.Panicln()
		return nil
	}
	return db
}
