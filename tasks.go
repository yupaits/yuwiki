package yuwiki

import (
	"github.com/robfig/cron"
	"log"
	"os"
	"path"
	"strings"
)

func StartScheduler() {
	backupCron := cron.New()
	if err := backupCron.AddFunc(Config.Cron.Backup, backup); err != nil {
		log.Println(err)
	} else {
		backupCron.Start()
	}
}

func backup() {
	dbFile := Config.DataSource.Url
	suffix := path.Ext(dbFile)
	dbFilename := path.Base(dbFile)
	prefix := Config.Path.BackupPath + dbFilename[0:strings.LastIndex(dbFilename, suffix)] + "."
	todayBackupFile := prefix + Today(DateLayout) + suffix
	yesterdayBackupFile := prefix + Yesterday(DateLayout) + suffix
	dbFileSha1, _ := FileSha1(dbFile)
	yesterdayFileSha1, _ := FileSha1(yesterdayBackupFile)
	if dbFileSha1 == yesterdayFileSha1 {
		//无增量数据时直接重命名备份文件
		if err := os.Rename(yesterdayBackupFile, todayBackupFile); err != nil {
			log.Printf("重命名备份文件失败，文件名：%s， 错误信息：%v", yesterdayBackupFile, err)
		} else {
			log.Printf("重命名备份文件成功，文件名：%s, 新文件名：%s", yesterdayBackupFile, todayBackupFile)
		}
	} else if size, err := CopyFile(todayBackupFile, dbFile); err != nil {
		log.Printf("数据库备份失败，错误信息：%v", err)
	} else {
		log.Printf("数据库备份成功，文件名：%s，大小：%s", todayBackupFile, ByteSize(uint64(size)))
	}
}
