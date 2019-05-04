package yuwiki

import (
	"github.com/robfig/cron"
	"log"
	"os"
	"strings"
)

func StartScheduler() {
	backupCron := cron.New()
	if err := backupCron.AddFunc(Config.Cron.Backup, backup); err != nil {
		log.Fatalln(err)
	} else {
		backupCron.Start()
	}
}

func backup() {
	dbFile := Config.DataSource.Url
	suffix := ".db"
	prefix := dbFile[0:strings.LastIndex(dbFile, suffix)] + "."
	todayBackupFile := prefix + Today(DateLayout) + suffix
	yesterdayBackupFile := prefix + Yesterday(DateLayout) + suffix
	dbFileSha1, _ := FileSha1(dbFile)
	yesterdayFileSha1, _ := FileSha1(yesterdayBackupFile)
	if dbFileSha1 == yesterdayFileSha1 {
		//无增量更新数据时直接重命名备份文件
		if err := os.Rename(yesterdayBackupFile, todayBackupFile); err != nil {
			log.Fatalf("重命名备份文件失败，文件名：%s", yesterdayBackupFile)
		}
	} else if size, err := CopyFile(todayBackupFile, dbFile); err != nil {
		log.Fatalf("数据库备份失败，错误信息：%v", err)
	} else {
		log.Printf("数据库备份成功，文件名：%s，大小：%s", todayBackupFile, ByteSize(uint64(size)))
	}
}
