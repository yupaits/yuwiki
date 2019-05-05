package yuwiki

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

type UploadFile struct {
	gorm.Model
	Filename     string `gorm:"unique;not null"`
	FilePath     string `gorm:"not null"`
	OriginalName string
	Size         int64
}

func uploadFileHandler(c *gin.Context) {
	if file, err := c.FormFile("file"); err != nil {
		log.Fatalf("文件上传失败，文件名：%s，错误原因：%v", file.Filename, err)
		Result(c, CodeFail(FileUploadFail))
	} else {
		//TODO
	}
}

func downloadFileHandler(c *gin.Context) {
	//TODO
}
