package yuwiki

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
	"time"
)

const (
	Secret = iota
	Male
	Female
)

const (
	PART = iota
	GROUP
)

type User struct {
	gorm.Model
	Username        string    `gorm:"unique;not null"`
	Password        string    `gorm:"size:60;not null"`
	Avatar          string    `json:"avatar"`
	Nickname        string    `json:"nickname"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Gender          int8      `json:"gender"`
	Birthday        time.Time `json:"birthday"`
	Admin           bool      `gorm:"not null" json:"admin"`
	Salt            string    `gorm:"not null"`
	InitPassword    string    `gorm:"size:60"`
	PasswordChanged bool
}

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
	Password  string `gorm:"size:60" json:"password"`
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

type TreePart struct {
	Part     Part
	SubParts *[]TreePart
}

func currentUser() (*User, error) {
	userId := GetUserId()
	user := &User{}
	if err := Db.Where("id = ?", userId).Find(user).Error; err != nil {
		log.Fatal("获取当前用户信息失败", err)
		return user, err
	}
	return user, nil
}

func GetCurrentUser() *User {
	user, _ := currentUser()
	user.Password = ""
	user.Salt = ""
	user.InitPassword = ""
	return user
}

func SaveUser(user *User) bool {
	var err error
	if Db.NewRecord(user) {
		//设置初始化密码
		salt := GenSalt()
		user.Salt = salt
		password := GenPassword()
		user.InitPassword = password
		user.PasswordChanged = false
		user.Password, _ = EncPassword(password, salt)
		err = Db.Create(user).Error
	} else {
		user.InitPassword = ""
		user.Password = ""
		user.Salt = ""
		err = Db.Model(user).Updates(user).Error
	}
	return err == nil
}

func modifyPassword(modify PasswordModify) (bool, string) {
	user, _ := currentUser()
	if modify.OldPassword == "" || modify.NewPassword == "" || modify.ConfirmPassword == "" {
		return false, "原密码、新密码、确认密码均不能为空"
	}
	if strings.Compare(modify.NewPassword, modify.ConfirmPassword) == 0 {
		return false, "新密码和确认密码不匹配"
	}
	if strings.Compare(modify.OldPassword, modify.NewPassword) == 0 {
		return false, "新密码不能与原密码相同"
	}
	if !Match(modify.OldPassword, user.Salt, user.Password) {
		return false, "原密码不正确"
	}

	password, err := EncPassword(modify.NewPassword, user.Salt)
	if err != nil {
		return false, err.Error()
	}
	user.Password = password
	user.PasswordChanged = true
	return true, ""
}

func getBooks() *[]Book {
	books := &[]Book{}
	owner := GetUserId()
	if err := Db.Where("owner = ?", owner).Find(books).Error; err != nil {
		log.Fatal("获取笔记本清单失败", err)
	}
	return books
}

func saveBook(book *Book) bool {
	var err error
	if Db.NewRecord(book) {
		err = Db.Create(book).Error
	} else {
		err = Db.Model(book).Updates(book).Error
	}
	return err == nil
}

func deleteBook(id uint) bool {
	var partCount, pageCount uint
	Db.Model(&Part{}).Where("book_id = ?", id).Count(&partCount)
	Db.Model(&Page{}).Where("book_id = ?", id).Count(&pageCount)
	if partCount+pageCount > 0 {
		return false
	}
	err := Db.Where("id = ?", id).Delete(Book{}).Error
	return err == nil
}

func getBookParts(bookId uint) *[]TreePart {
	parts := &[]Part{}
	if err := Db.Where("book_id = ? AND parent_Id = 0", bookId).Find(parts).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取笔记本分区清单失败，bookId: %d", bookId), err)
	}
	var treeParts []TreePart
	for _, part := range *parts {
		if part.Protected && part.Password != "" {
			part.Password = ""
		}
		treePart := TreePart{
			Part:     part,
			SubParts: getSubParts(part.ID),
		}
		treeParts = append(treeParts, treePart)
	}
	return &treeParts
}

func getSubParts(parentId uint) *[]TreePart {
	parts := &[]Part{}
	if err := Db.Where("parent_id = ?", parentId).Find(parts).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取笔记本分区子分区列表失败，parentId: %d", parentId), err)
	}
	var subParts []TreePart
	for _, part := range *parts {
		if part.Protected && part.Password != "" {
			part.Password = ""
		}
		subParts = append(subParts, TreePart{
			Part:     part,
			SubParts: getSubParts(part.ID),
		})
	}
	return &subParts
}

func savePart(part *Part) (bool, error) {
	if part.Protected {
		user, err := currentUser()
		if err != nil {
			return false, err
		}
		if part.Password, err = EncPassword(part.Password, user.Salt); err != nil {
			return false, err
		}
	}
	var err error
	if Db.NewRecord(part) {
		err = Db.Create(part).Error
	} else {
		err = Db.Model(part).Updates(part).Error
	}
	return err == nil, err
}

func deletePart(id uint) bool {
	var subPartCount, pageCount uint
	Db.Model(&Part{}).Where("parent_id = ?", id).Count(&subPartCount)
	Db.Model(&Page{}).Where("part_id = ?", id).Count(&pageCount)
	if subPartCount+pageCount > 0 {
		return false
	}
	err := Db.Where("id = ?", id).Delete(Part{}).Error
	return err == nil
}

func getPartPages(partId uint) *[]Page {
	pages := &[]Page{}
	if err := Db.Where("part_id = ?", partId).Find(pages).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取分区页面清单失败，partId: %d", partId), err)
	}
	for _, page := range *pages {
		page.Content = ""
	}
	return pages
}

func getPage(id uint) *Page {
	page := &Page{}
	if err := Db.Where("id = ?", id).Find(page).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取页面失败，pageId: %d", id), err)
	}
	return page
}

func savePage(page *Page) bool {
	var err error
	if Db.NewRecord(page) {
		err = Db.Create(page).Error
	} else {
		err = Db.Model(page).Updates(page).Error
	}
	return err == nil
}

func deletePage(id uint) bool {
	err := Db.Where("id = ?", id).Delete(Page{}).Error
	return err == nil
}
