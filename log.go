package yuwiki

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

var log *logrus.Logger

func InitLog() {
	log = logrus.New()
	Mkdirs(Config.LogFile)
	//logFile, err := os.OpenFile(Config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	//切分日志文件
	initLfsHook()
	//控制台日志
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: DateTimeLogLayout})
}

func initLfsHook() {
	logFile := Config.LogFile
	suffix := path.Ext(logFile)
	prefix := logFile[0:strings.LastIndex(logFile, suffix)]
	writer, err := rotatelogs.New(
		prefix+".%Y%m%d"+suffix,
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationCount(Config.LogFileMaxCount),
	)
	if err != nil {
		log.WithField("error", err).Error("日志文件切分配置出错")
	}
	var formatter logrus.Formatter
	if Config.Debug {
		formatter = &logrus.TextFormatter{DisableColors: true, TimestampFormat: DateTimeLogLayout}
		log.SetLevel(logrus.DebugLevel)
	} else {
		formatter = &logrus.JSONFormatter{TimestampFormat: DateTimeLogLayout}
		log.SetLevel(logrus.InfoLevel)
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatter)
	log.AddHook(lfsHook)
}
