package yuwiki

import (
	"errors"
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
	SortCode  uint   `gorm:"not null"`
}

type Page struct {
	gorm.Model
	BookId    uint   `gorm:"not null;index" json:"bookId" binding:"required"`
	PartId    uint   `gorm:"not null;index" json:"partId" binding:"required"`
	Title     string `gorm:"not null;index" json:"title" binding:"required"`
	Content   string `gorm:"type:text" json:"content"`
	Tags      *[]Tag `gorm:"many2many:page_tags" json:"tags"`
	Published bool
	Owner     uint `gorm:"not null"`
}

type Tag struct {
	gorm.Model
	Name  string  `gorm:"not null;unique" json:"name" binding:"required"`
	Pages *[]Page `gorm:"many2many:page_tags"`
}

type HistoricalPage struct {
	ID        uint   `gorm:"primary_key"`
	PageId    uint   `gorm:"not null;index"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}

type TreePart struct {
	Part
	SubParts *[]TreePart
}

type SharedBook struct {
	gorm.Model
	BookId uint `gorm:"not null;index" json:"bookId" binding:"required"`
	UserId uint `gorm:"not null;index" json:"userId" binding:"required"`
}

func currentUser() (*User, error) {
	userId := getUserId()
	user := &User{}
	if err := Db.Where("id = ?", userId).Find(user).Error; err != nil {
		log.Fatal("获取当前用户信息失败", err)
		return user, err
	}
	return user, nil
}

func getCurrentUser() (*User, error) {
	if user, err := currentUser(); err != nil {
		return nil, err
	} else {
		user.Password = ""
		user.Salt = ""
		user.InitPassword = ""
		return user, nil
	}
}

func saveUser(user *User) bool {
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
		err = Db.Save(user).Error
	}
	return err == nil
}

func modifyPassword(modify *PasswordModify) (bool, string) {
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
	owner := getUserId()
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
		err = Db.Save(book).Error
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
	err := Db.Where("id = ? AND owner = ?", id, getUserId()).Delete(Book{}).Error
	return err == nil
}

func getBookParts(bookId uint) *[]TreePart {
	parts := &[]Part{}
	if err := Db.Where("book_id = ? AND parent_Id = 0 AND owner = ?", bookId, getUserId()).Find(parts).Error; err != nil {
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
		subParts = append(subParts, TreePart{part, getSubParts(part.ID)})
	}
	if subParts == nil {
		return &[]TreePart{}
	} else {
		return &subParts
	}
}

func getPart(partId uint) *Part {
	part := &Part{}
	if err := Db.Where("id = ? AND owner = ?", partId, getUserId()).Find(part).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取分区信息失败，partId: %d", partId), err)
	}
	return part
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
	var max uint
	if rows, err := Db.Table("parts").Select(" MAX(parts.sortCode) AS max").Where("parent_id = ?", part.ParentId).Rows(); err == nil {
		if rows.Next() {
			if err := rows.Scan(&max); err != nil {
				log.Fatal("取排序码最大值出错", err)
			}
		}
	}
	part.SortCode = max + 1
	var err error
	if Db.NewRecord(part) {
		err = Db.Create(part).Error
	} else {
		dbPart := &Part{}
		if err := Db.Where("id = ? AND owner = ?", part.ID, getUserId()).Find(dbPart).Error; err != nil {
			return false, err
		} else if dbPart.PartType == GROUP && part.PartType == PART {
			var subPartCount uint
			Db.Model(&Part{}).Where("parent_id = ?", part.ID).Count(&subPartCount)
			if subPartCount > 0 {
				return false, errors.New("包含子分区的分区组不能修改分区类型")
			}
		} else if dbPart.PartType == PART && part.PartType == GROUP {
			var pageCount uint
			Db.Model(&Page{}).Where("part_id = ? AND owner = ?", part.ID, getUserId()).Count(&pageCount)
			if pageCount > 0 {
				return false, errors.New("包含页面的分区不能修改分区类型")
			}
		}
		err = Db.Save(part).Error
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
	err := Db.Where("id = ? AND owner = ?", id, getUserId()).Delete(Part{}).Error
	return err == nil
}

func getPartPages(partId uint) *[]Page {
	pages := &[]Page{}
	if err := Db.Where("part_id = ? AND owner = ?", partId, getUserId()).Find(pages).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取分区页面清单失败，partId: %d", partId), err)
	}
	for _, page := range *pages {
		page.Content = ""
	}
	return pages
}

func getPage(id uint) *Page {
	page := &Page{}
	if err := Db.Where("id = ? AND owner = ?", id, getUserId()).Find(page).Error; err != nil {
		log.Fatal(fmt.Sprintf("获取页面失败，pageId: %d", id), err)
	}
	//页面处于草稿状态时，返回最近发布的页面内容
	if !page.Published {
		historicalPage := &HistoricalPage{}
		if err := Db.Where("page_id = ?").Order("created_at DESC").Limit(1).Find(historicalPage); err == nil {
			page.Content = historicalPage.Content
		}
	}
	return page
}

func savePage(page *Page) bool {
	var err error
	if Db.NewRecord(page) {
		err = Db.Create(page).Error
	} else {
		err = Db.Save(page).Error
	}
	return err == nil
}

func editPage(page *Page) bool {
	if savePage(page) {
		//发布状态页面需要保存页面历史记录
		if page.Published {
			historicalPage := &HistoricalPage{
				PageId:    page.ID,
				Content:   page.Content,
				CreatedAt: time.Now(),
			}
			if err := Db.Create(historicalPage).Error; err != nil {
				log.Fatal(fmt.Sprintf("保存页面历史记录失败, pageId: %d", page.ID), err)
			}
		}
		return true
	} else {
		return false
	}
}

func deletePage(id uint) bool {
	err := Db.Where("id = ? AND owner = ?", id, getUserId()).Delete(Page{}).Error
	return err == nil
}

func getSharedBooks() *[]Book {
	sharedBooks := &[]SharedBook{}
	Db.Where("user_id = ?", getUserId()).Find(sharedBooks)
	var bookIds []uint
	for _, sharedBook := range *sharedBooks {
		bookIds = append(bookIds, sharedBook.BookId)
	}
	books := &[]Book{}
	Db.Where("book_id in (?)", bookIds).Find(books)
	return books
}

func saveSharedBook(sharedBook *SharedBook) bool {
	var err error
	if Db.NewRecord(sharedBook) {
		err = Db.Create(sharedBook).Error
	} else {
		err = Db.Save(sharedBook).Error
	}
	return err == nil
}

func getStarItems() *StarItems {
	var starBooks []Book
	var starParts []Part
	var starPages []Page
	starItems := &StarItems{
		Books: &starBooks,
		Parts: &starParts,
		Pages: &starPages,
	}
	return starItems
}
