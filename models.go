package yuwiki

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	Secret = iota + 1
	Male
	Female
)

const (
	PART = iota
	GROUP
)

type User struct {
	gorm.Model
	Username        string    `gorm:"unique;not null" json:"username"`
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
	Name     string `gorm:"unique;not null" json:"name" binding:"required"`
	Color    string `json:"color"`
	Star     bool   `json:"star"`
	Owner    uint   `gorm:"not null" json:"owner"`
	SortCode uint   `gorm:"not null" json:"sortCode"`
}

type Part struct {
	gorm.Model
	BookId    uint   `gorm:"not null" json:"bookId" binding:"required"`
	ParentId  uint   `gorm:"not null" json:"parentId"`
	Name      string `gorm:"not null" json:"name" binding:"required"`
	PartType  int8   `gorm:"not null" json:"partType"`
	Protected bool   `json:"protected"`
	Password  string `gorm:"size:60" json:"password"`
	Star      bool   `json:"star"`
	Owner     uint   `gorm:"not null" json:"owner"`
	SortCode  uint   `gorm:"not null" json:"sortCode"`
}

type Page struct {
	gorm.Model
	BookId    uint   `gorm:"not null;index" json:"bookId"`
	PartId    uint   `gorm:"not null;index" json:"partId"`
	Title     string `gorm:"not null;index" json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	Published bool   `json:"published"`
	Star      bool   `json:"star"`
	Owner     uint   `gorm:"not null" json:"owner"`
	SortCode  uint   `gorm:"not null" json:"sortCode"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"`
}

type PageTag struct {
	gorm.Model
	PageId uint `gorm:"not null;index"`
	TagId  uint `gorm:"not null;index"`
}

type HistoricalPage struct {
	ID        uint   `gorm:"primary_key"`
	PageId    uint   `gorm:"not null;index" json:"page_id"`
	Content   string `gorm:"type:text" json:"content"`
	CreatedAt time.Time
	Owner     uint `gorm:"not null" json:"owner"`
}

type TreePart struct {
	Part
	SubParts *[]TreePart
}

type PageVo struct {
	Page
	Tags []string `json:"tags"`
}

type SharedUser struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type SharedBook struct {
	gorm.Model
	BookId uint `gorm:"not null;index" json:"bookId" binding:"required"`
	UserId uint `gorm:"not null;index" json:"userId" binding:"required"`
}

type maxSortCode struct {
	Max uint
}

func currentUser() (*User, error) {
	userId := getUserId()
	user := &User{}
	if err := Db.Where("id = ?", userId).Find(user).Error; err != nil {
		log.WithField("error", err).Error("获取当前用户信息失败")
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
		if user.Password == "" {
			password := GenPassword()
			user.InitPassword = password
			user.PasswordChanged = false
			user.Password, _ = EncPassword(password, salt)
		}
		err = Db.Create(user).Error
	} else {
		user.InitPassword = ""
		user.Password = ""
		user.Salt = ""
		err = Db.Save(user).Error
	}
	return err == nil
}

func searchUsers(userSearch *UserSearch) *[]User {
	users := &[]User{}
	var query string
	if userSearch.ID > 0 {
		query = fmt.Sprintf("id = %d", userSearch.ID)
	}
	if userSearch.Gender > 0 {
		genderQuery := fmt.Sprintf("gender = %d", userSearch.Gender)
		if query == "" {
			query = genderQuery
		} else {
			query = query + " AND " + genderQuery
		}
	}
	if userSearch.NoSelf {
		noSelfQuery := fmt.Sprintf("id != %d", getUserId())
		if query == "" {
			query = noSelfQuery
		} else {
			query = query + " AND " + noSelfQuery
		}
	}
	if userSearch.Keyword != "" {
		keyword := "'%" + userSearch.Keyword + "%'"
		keywordQuery := fmt.Sprintf("username LIKE %s OR nickname LIKE %s OR phone LIKE %s OR email LIKE %s", keyword, keyword, keyword, keyword)
		if query == "" {
			query = keywordQuery
		} else {
			query = query + " AND (" + keywordQuery + ")"
		}
	}
	if err := Db.Where(query).Find(users).Error; err != nil {
		log.WithField("error", err).Error("查询用户失败")
	}
	var resUsers []User
	for _, user := range *users {
		user.InitPassword = ""
		user.Password = ""
		user.Salt = ""
		resUsers = append(resUsers, user)
	}
	return &resUsers
}

func modifyPassword(modify *PasswordModify) (bool, string) {
	user, _ := currentUser()
	if modify.OldPassword == "" || modify.NewPassword == "" || modify.ConfirmPassword == "" {
		return false, "原密码、新密码、确认密码均不能为空"
	}
	if strings.Compare(modify.NewPassword, modify.ConfirmPassword) != 0 {
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
	user.InitPassword = ""
	if err := Db.Save(user).Error; err != nil {
		return false, "更新密码失败"
	}
	return true, ""
}

func getBooks() *[]Book {
	books := &[]Book{}
	if err := Db.Where("owner = ?", getUserId()).Order("sort_code").Find(books).Error; err != nil {
		log.WithField("error", err).Error("获取笔记本清单失败")
	}
	return books
}

func saveBook(book *Book) bool {
	if book.SortCode == 0 {
		var maxCode maxSortCode
		if err := Db.Raw("SELECT MAX(books.sort_code) AS max FROM books WHERE owner = ?", getUserId()).Scan(&maxCode).Error; err != nil {
			maxCode.Max = 0
		}
		book.SortCode = maxCode.Max + 1
	}
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
	if err := Db.Where("book_id = ? AND parent_Id = 0 AND owner = ?", bookId, getUserId()).Order("sort_code").Find(parts).Error; err != nil {
		log.WithFields(logrus.Fields{
			"bookId": bookId,
			"error":  err,
		}).Error("获取笔记本分区清单失败")
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
	if err := Db.Where("parent_id = ?", parentId).Order("sort_code").Find(parts).Error; err != nil {
		log.WithFields(logrus.Fields{
			"parentId": parentId,
			"error":    err,
		}).Error("获取笔记本分区子分区列表失败")
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
		log.WithFields(logrus.Fields{
			"partId": partId,
			"error":  err,
		}).Error("获取分区信息失败")
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
	if part.SortCode == 0 {
		var maxCode maxSortCode
		if err := Db.Raw("SELECT MAX(parts.sort_code) AS max FROM parts WHERE parent_id = ? AND owner = ?",
			part.ParentId, getUserId()).Scan(&maxCode).Error; err != nil {
			maxCode.Max = 0
		}
		part.SortCode = maxCode.Max + 1
	}
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

func getPartPages(partId uint) *[]PageVo {
	pages := &[]Page{}
	if err := Db.Where("part_id = ? AND owner = ?", partId, getUserId()).Order("sort_code").Find(pages).Error; err != nil {
		log.WithFields(logrus.Fields{
			"partId": partId,
			"error":  err,
		}).Error("获取分区页面清单失败")
	}
	var pageVos []PageVo
	for _, page := range *pages {
		pageVo := PageVo{Page: page}
		pageVo.Content = ""
		setPageTags(&pageVo)
		pageVos = append(pageVos, pageVo)
	}
	return &pageVos
}

func getPage(id uint, editable bool) *PageVo {
	page := &Page{}
	if err := Db.Where("id = ? AND owner = ?", id, getUserId()).Find(page).Error; err != nil {
		log.WithFields(logrus.Fields{
			"pageId": id,
			"error":  err,
		}).Error("获取页面失败")
	}
	pageVo := PageVo{Page: *page}
	setPageTags(&pageVo)
	//页面处于草稿状态时，返回最近发布的页面内容
	if !page.Published && !editable {
		historicalPage := &HistoricalPage{}
		Db.Where("page_id = ?", page.ID).Order("created_at DESC").Limit(1).Find(historicalPage)
		pageVo.Content = historicalPage.Content
	}
	return &pageVo
}

//填充页面标签信息
func setPageTags(pageVo *PageVo) {
	var pageTags []PageTag
	if err := Db.Where("page_id = ?", pageVo.ID).Find(&pageTags).Error; err != nil {
		log.WithFields(logrus.Fields{
			"pageId": pageVo.ID,
			"error":  err,
		}).Error("获取页面标签列表失败")
	} else if len(pageTags) > 0 {
		var tagIds []uint
		for _, pageTag := range pageTags {
			tagIds = append(tagIds, pageTag.TagId)
		}
		tags := &[]Tag{}
		if err := Db.Where("id in (?)", tagIds).Find(tags).Error; err != nil {
			log.WithFields(logrus.Fields{
				"tagIds": tagIds,
				"error":  err,
			}).Error("获取标签信息失败")
		}
		pageVo.Tags = []string{}
		for _, tag := range *tags {
			pageVo.Tags = append(pageVo.Tags, tag.Name)
		}
	}
}

func savePage(pageDto *PageDto) bool {
	page := &Page{
		Model:     gorm.Model{ID: pageDto.ID, CreatedAt: pageDto.CreatedAt},
		BookId:    pageDto.BookId,
		PartId:    pageDto.PartId,
		Title:     pageDto.Title,
		Content:   pageDto.Content,
		Published: pageDto.Published,
		Owner:     pageDto.Owner,
		SortCode:  pageDto.SortCode,
	}
	if page.SortCode == 0 {
		var maxCode maxSortCode
		if err := Db.Raw("SELECT MAX(pages.sort_code) AS max FROM pages WHERE owner = ?", getUserId()).Scan(&maxCode).Error; err != nil {
			maxCode.Max = 0
		}
		page.SortCode = maxCode.Max + 1
	}
	var err error
	if Db.NewRecord(page) {
		err = Db.Create(page).Error
	} else {
		err = Db.Save(page).Error
	}
	var tagErr error
	if len(pageDto.Tags) > 0 {
		//更新页面标签
		var tagIds []uint
		for _, tagName := range pageDto.Tags {
			tag := &Tag{}
			Db.Where("name = ?", tagName).First(tag)
			if tag.ID == 0 {
				//创建不存在的标签
				tag.Name = tagName
				Db.Create(tag)
			}
			tagIds = append(tagIds, tag.ID)
			pageTag := &PageTag{}
			Db.Where("page_id = ? AND tag_id = ?", page.ID, tag.ID).First(pageTag)
			if pageTag.ID == 0 {
				pageTag.PageId = page.ID
				pageTag.TagId = tag.ID
				//添加新的页面标签关联记录
				if err := Db.Create(pageTag).Error; err != nil {
					log.WithFields(logrus.Fields{
						"pageId": page.ID,
						"tagId":  tag.ID,
						"error":  err,
					}).Error("创建页面标签关联记录失败")
				}
			}
		}
		//删除无效的页面标签关联记录
		if err := Db.Where("page_id = ?", page.ID).Not("tag_id", tagIds).
			Delete(PageTag{}).Error; err != nil {
			log.WithFields(logrus.Fields{
				"pageId":      page.ID,
				"notInTagIds": tagIds,
				"error":       err,
			}).Error("删除页面标签管理记录失败")
		}
	} else {
		tagErr = Db.Where("page_id = ?", page.ID).Delete(PageTag{}).Error
	}
	if tagErr != nil {
		log.WithFields(logrus.Fields{
			"pageId": page.ID,
			"error":  err,
		}).Error("更新页面标签信息失败")
	}
	return err == nil
}

func editPage(pageDto *PageDto) bool {
	if savePage(pageDto) {
		//发布状态页面需要保存页面历史记录
		if pageDto.Published {
			historicalPage := &HistoricalPage{
				PageId:    pageDto.ID,
				Content:   pageDto.Content,
				CreatedAt: time.Now(),
				Owner:     getUserId(),
			}
			if err := Db.Create(historicalPage).Error; err != nil {
				log.WithFields(logrus.Fields{
					"pageId": pageDto.ID,
					"error":  err,
				}).Error("保存页面历史记录失败")
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

func getTags() *[]Tag {
	tags := &[]Tag{}
	if err := Db.Find(tags).Error; err != nil {
		log.WithField("error", err).Error("获取页面标签列表失败")
	}
	return tags
}

func getHistoricalPages(pageId uint) *[]HistoricalPage {
	historicalPages := &[]HistoricalPage{}
	if err := Db.Where("page_id = ? AND owner = ?", pageId, getUserId()).Order("created_at DESC").Find(historicalPages).Error; err != nil {
		log.WithFields(logrus.Fields{
			"pageId": pageId,
			"error":  err,
		}).Error("获取页面历史记录失败")
	}
	return historicalPages
}

func saveSharedBook(sharedBook *SharedBook) bool {
	var err error
	dbSharedBook := &SharedBook{}
	Db.Where("user_id = ? AND book_id = ?", sharedBook.UserId, sharedBook.BookId).Limit(1).Find(dbSharedBook)
	if dbSharedBook.ID == 0 {
		err = Db.Create(sharedBook).Error
	}
	return err == nil
}

func deleteSharedBook(sharedBook *SharedBook) bool {
	err := Db.Where("user_id = ? AND book_id = ?", sharedBook.UserId, sharedBook.BookId).Delete(SharedBook{}).Error
	return err == nil
}

func getBookSharedUsers(bookId uint) *[]SharedUser {
	sharedBooks := &[]SharedBook{}
	Db.Where("book_id =?", bookId).Find(sharedBooks)
	var userIds []uint
	for _, sharedBook := range *sharedBooks {
		userIds = append(userIds, sharedBook.UserId)
	}
	users := &[]User{}
	Db.Where("id in (?)", userIds).Find(users)
	var sharedUsers []SharedUser
	for _, user := range *users {
		sharedUser := SharedUser{
			UserId:   user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		}
		sharedUsers = append(sharedUsers, sharedUser)
	}
	return &sharedUsers
}

func searchPagesByKeyword(keyword string) *[]PageVo {
	keyword = "%" + keyword + "%"
	pages := &[]Page{}
	if err := Db.Where("(title LIKE ? OR content LIKE ?) AND owner = ?", keyword, keyword, getUserId()).Order("updated_at DESC").Find(pages).Error; err != nil {
		log.WithFields(logrus.Fields{
			"keyword": keyword,
			"error":   err,
		}).Error("根据关键字查找页面失败")
	}
	var pageVos []PageVo
	for _, page := range *pages {
		pageVo := PageVo{Page: page}
		pageVo.Content = ""
		setPageTags(&pageVo)
		pageVos = append(pageVos, pageVo)
	}
	return &pageVos
}

func sortBooks(sortedBooks *[]SortBook) (bool, error) {
	tx := Db.Begin()
	for _, sortedBook := range *sortedBooks {
		book := Book{Model: gorm.Model{ID: sortedBook.BookId}}
		if err := tx.Model(&book).Update("sort_code", sortedBook.SortCode).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}
	tx.Commit()
	return true, nil
}

func sortParts(sortedParts *[]SortPart) (bool, error) {
	tx := Db.Begin()
	for _, sortedPart := range *sortedParts {
		part := Part{Model: gorm.Model{ID: sortedPart.PartId}}
		if err := tx.Model(&part).Update("sort_code", sortedPart.SortCode).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}
	tx.Commit()
	return true, nil
}

func sortPages(sortedPages *[]SortPage) (bool, error) {
	tx := Db.Begin()
	for _, sortedPage := range *sortedPages {
		page := Page{Model: gorm.Model{ID: sortedPage.PageId}}
		if err := tx.Model(&page).Update("sort_code", sortedPage.SortCode).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}
	tx.Commit()
	return true, nil
}

func toggleStarBook(bookId uint) bool {
	book := &Book{}
	if err := Db.Where("id = ? AND owner = ?", bookId, getUserId()).Find(book).Error; err != nil {
		return false
	} else if book.ID == 0 {
		return false
	} else {
		book.Star = !book.Star
		if err := Db.Save(book).Error; err != nil {
			return false
		}
		return true
	}
}

func toggleStarPart(partId uint) bool {
	part := &Part{}
	if err := Db.Where("id = ? AND owner = ?", partId, getUserId()).Find(part).Error; err != nil {
		return false
	} else if part.ID == 0 {
		return false
	} else {
		part.Star = !part.Star
		if err := Db.Save(part).Error; err != nil {
			return false
		}
		return true
	}
}

func toggleStarPage(pageId uint) bool {
	page := &Page{}
	if err := Db.Where("id = ? AND owner = ?", pageId, getUserId()).Find(page).Error; err != nil {
		return false
	} else if page.ID == 0 {
		return false
	} else {
		page.Star = !page.Star
		if err := Db.Save(page).Error; err != nil {
			return false
		}
		return true
	}
}

func getSharedBooks() *[]Book {
	sharedBooks := &[]SharedBook{}
	Db.Where("user_id = ?", getUserId()).Find(sharedBooks)
	var bookIds []uint
	for _, sharedBook := range *sharedBooks {
		bookIds = append(bookIds, sharedBook.BookId)
	}
	books := &[]Book{}
	Db.Where("id in (?)", bookIds).Find(books)
	return books
}

func getSharedParts(bookId uint) *[]TreePart {
	var treeParts []TreePart
	sharedBook := &SharedBook{}
	Db.Where("book_id = ? AND user_id = ?", bookId, getUserId()).Find(sharedBook)
	if sharedBook.ID != 0 {
		parts := &[]Part{}
		Db.Where("book_id = ?", bookId).Find(parts)
		if err := Db.Where("book_id = ? AND parent_Id = 0", bookId).Order("sort_code").Find(parts).Error; err != nil {
			log.WithFields(logrus.Fields{
				"bookId": bookId,
				"error":  err,
			}).Error("获取笔记本分区清单失败")
		}
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
	}
	return &treeParts
}

func getSharedPages(bookId, partId uint) *[]PageVo {
	var pageVos []PageVo
	sharedBook := &SharedBook{}
	Db.Where("book_id = ? AND user_id = ?", bookId, getUserId()).Find(sharedBook)
	if sharedBook.ID != 0 {
		pages := &[]Page{}
		if err := Db.Where("book_id = ? AND part_id = ?", bookId, partId).Order("sort_code").Find(pages).Error; err != nil {
			log.WithFields(logrus.Fields{
				"bookId": bookId,
				"partId": partId,
				"error":  err,
			}).Error("获取分区页面清单失败")
		}
		for _, page := range *pages {
			pageVo := PageVo{Page: page}
			pageVo.Content = ""
			setPageTags(&pageVo)
			pageVos = append(pageVos, pageVo)
		}
	}
	return &pageVos
}

func getSharedPage(bookId, partId, pageId uint) *PageVo {
	var pageVo PageVo
	sharedBook := &SharedBook{}
	Db.Where("book_id = ? AND user_id = ?", bookId, getUserId()).Find(sharedBook)
	if sharedBook.ID != 0 {
		page := &Page{}
		if err := Db.Where("id = ? AND book_id = ? AND part_id = ?", pageId, bookId, partId).Find(page).Error; err != nil {
			log.WithFields(logrus.Fields{
				"bookId": bookId,
				"partId": partId,
				"pageId": pageId,
				"error":  err,
			}).Error("获取页面失败")
		}
		pageVo = PageVo{Page: *page}
		setPageTags(&pageVo)
		//页面处于草稿状态时，返回最近发布的页面内容
		if !page.Published {
			historicalPage := &HistoricalPage{}
			Db.Where("page_id = ?", page.ID).Order("created_at DESC").Limit(1).Find(historicalPage)
			pageVo.Content = historicalPage.Content
		}
	}
	return &pageVo
}
