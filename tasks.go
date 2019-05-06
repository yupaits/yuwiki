package yuwiki

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

func StartScheduler() {
	backupCron := cron.New()
	if err := backupCron.AddFunc(Config.Cron.Backup, backup); err != nil {
		log.Error(err)
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
			log.WithFields(logrus.Fields{
				"filename": yesterdayBackupFile,
				"error":    err,
			}).Error("重命名备份文件失败")
		} else {
			log.WithFields(logrus.Fields{
				"filename": yesterdayBackupFile,
				"newFile":  todayBackupFile,
			}).Info("重命名备份文件成功")
		}
	} else if size, err := CopyFile(todayBackupFile, dbFile); err != nil {
		log.WithField("error", err).Error("数据库备份失败")
	} else {
		log.WithFields(logrus.Fields{
			"filename": todayBackupFile,
			"size":     ByteSize(uint64(size)),
		}).Info("数据库备份成功")
	}
}
