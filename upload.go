package yuwiki

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"path"
)

type UploadFile struct {
	gorm.Model
	Filename     string `gorm:"unique;not null"`
	FilePath     string `gorm:"not null"`
	OriginalName string
	Size         int64
}

func uploadFileHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Errorf("文件上传失败，错误原因：%v", err)
		Result(c, CodeFail(FileUploadFail))
		return
	}
	if fileHeader == nil {
		log.Error("上传文件为空")
		Result(c, CodeFail(FileUploadFail))
		return
	}
	filename := NowFilename(DateTimeFileLayout) + GenFilename() + path.Ext(fileHeader.Filename)
	filepath := Config.Path.UploadPath + filename
	if err := c.SaveUploadedFile(fileHeader, filepath); err != nil {
		log.Errorf("保存文件失败，错误原因：%v", err)
		Result(c, CodeFail(FileUploadFail))
		return
	}
	uploadFile := &UploadFile{
		Filename:     filename,
		FilePath:     filepath,
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
	}
	if err := Db.Create(uploadFile).Error; err != nil {
		log.Errorf("保存文件上传记录失败，错误原因：%v", err)
		Result(c, CodeFail(FileUploadFail))
	} else {
		fileUrl := "/file/" + filename
		Result(c, OkData(fileUrl))
		log.Infof("上传文件成功，文件名：%s，保存路径：%s", filename, filepath)
	}
}

func downloadFileHandler(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		Result(c, CodeFail(ParamsError))
		return
	}
	file := &UploadFile{}
	if err := Db.Where("filename = ?", filename).Find(file).Error; err != nil {
		Result(c, MsgFail(fmt.Sprintf("获取文件%s上传记录失败", filename)))
		return
	}
	c.File(file.FilePath)
}
