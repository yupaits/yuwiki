package yuwiki

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

const (
	PART = iota
	GROUP
)

type Book struct {
	gorm.Model
	Name  string `gorm:"unique;not null" json:"name" binding:"required"`
	Color string `json:"color"`
	Star  bool   `json:"star"`
	Owner uint   `gorm:"not null"`
}

type Part struct {
	gorm.Model
	BookId    uint   `gorm:"not null" json:"bookId" binding:"required"`
	ParentId  uint   `gorm:"not null" json:"parentId"`
	Name      string `gorm:"not null" json:"name" binding:"required"`
	PartType  int8   `gorm:"not null" json:"partType"`
	Protected bool   `json:"protected"`
	Password  string `json:"password"`
	Owner     uint   `gorm:"not null"`
}

type Page struct {
	gorm.Model
	BookId  uint   `gorm:"not null;index" json:"bookId" binding:"required"`
	PartId  uint   `gorm:"not null;index" json:"partId" binding:"required"`
	Title   string `gorm:"not null;index" json:"title" binding:"required"`
	Content string `gorm:"type:text" json:"content"`
	Tags    *[]Tag `gorm:"many2many:page_tags" json:"tags"`
	Owner   uint   `gorm:"not null"`
}

type Tag struct {
	gorm.Model
	Name  string  `gorm:"not null;unique" json:"name" binding:"required"`
	Pages *[]Page `gorm:"many2many:page_tags"`
}

func GetBooks() *[]Book {
	db := DbConn()
	defer db.Close()
	books := &[]Book{}
	owner := GetUserId()
	if err := db.Where("Owner = ?", owner).Find(books).Error; err != nil {
		log.Fatal("获取笔记本清单失败", err)
	}
	return books
}

func SaveBook(book *Book) bool {
	db := DbConn()
	defer db.Close()
	var err error
	if db.NewRecord(book) {
		err = db.Create(book).Error
	} else {
		err = db.Updates(book).Error
	}
	return err == nil
}

func DeleteBook(id uint) bool {
	db := DbConn()
	defer db.Close()
	var partCount, pageCount uint
	db.Model(&Part{}).Where("BookId = ?", id).Count(&partCount)
	db.Model(&Page{}).Where("BookId = ?", id).Count(&pageCount)
	if partCount+pageCount > 0 {
		return false
	}
	return true
}

func GetBookParts(bookId uint) *[]Part {
	db := DbConn()
	defer db.Close()
	parts := &[]Part{}
	if err := db.Where("BookId = ? AND ParentId = 0", bookId).Find(parts).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取笔记本分区清单失败，bookId: %d", bookId), err)
	}
	for _, part := range *parts {
		if part.Protected && part.Password != "" {
			part.Password = ""
		}
	}
	return parts
}

func SavePart(part *Part) bool {
	db := DbConn()
	defer db.Close()
	//TODO 密码加密
	return true
}
