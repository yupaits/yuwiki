package yuwiki

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PasswordModify struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type StarItems struct {
	Books *[]Book
	Parts *[]Part
	Pages *[]Page
}

func getBooksHandler(c *gin.Context) {
	Result(c, OkData(getBooks()))
}

func getBookPartsHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getBookParts(uint(bookId))))
	}
}

func getPartHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPart(uint(partId))))
	}
}

func getSharedBooksHandler(c *gin.Context) {
	Result(c, OkData(getSharedBooks()))
}

func saveBookHandler(c *gin.Context) {
	book := &Book{}
	if err := c.ShouldBindJSON(book); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if saveBook(book) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func deleteBookHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deleteBook(uint(bookId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func shareBookHandler(c *gin.Context) {
	sharedBook := &SharedBook{}
	if err := c.ShouldBind(sharedBook); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if saveSharedBook(sharedBook) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func getPartPagesHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPartPages(uint(partId))))
	}
}

func savePartHandler(c *gin.Context) {
	part := &Part{}
	if err := c.ShouldBind(part); err != nil {
		Result(c, CodeFail(ParamsError))
		return
	}
	if part.Protected && part.Password == "" {
		Result(c, CodeFail(ParamsError))
		return
	}
	if ok, err := savePart(part); ok {
		Result(c, Ok())
	} else {
		Result(c, MsgFail(err.Error()))
	}
}

func deletePartHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deletePart(uint(partId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getPageHandler(c *gin.Context) {
	pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32)
	if err != nil {
		Result(c, CodeFail(ParamsError))
	}
	Result(c, OkData(getPage(uint(pageId))))
}

func savePageHandler(c *gin.Context) {
	page := &Page{}
	if err := c.ShouldBind(page); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if savePage(page) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func deletePageHandler(c *gin.Context) {
	if pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deletePage(uint(pageId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getUserInfoHandler(c *gin.Context) {
	if user, err := getCurrentUser(); err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, OkData(user))
	}
}

func modifyPasswordHandler(c *gin.Context) {
	modify := &PasswordModify{}
	if err := c.ShouldBind(modify); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, msg := modifyPassword(modify); ok {
		Result(c, Ok())
	} else {
		Result(c, MsgFail(msg))
	}
}

func getStarItemsHandler(c *gin.Context) {
	Result(c, OkData(getStarItems()))
}

func siteSearchHandler(c *gin.Context) {

}
