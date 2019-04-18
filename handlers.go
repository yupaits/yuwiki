package yuwiki

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type PageDto struct {
	ID        uint     `json:"id"`
	BookId    uint     `json:"bookId" binding:"required"`
	PartId    uint     `json:"partId" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Owner     uint     `json:"owner"`
	SortCode  uint     `json:"sortCode"`
	Published bool
	CreatedAt time.Time `json:"CreatedAt"`
}

type PasswordModify struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type StarItems struct {
	Books *[]Book `json:"books"`
	Parts *[]Part `json:"parts"`
	Pages *[]Page `json:"pages"`
}

type SortBook struct {
	BookId   uint `json:"bookId"`
	SortCode uint `json:"sortCode"`
}

type SortPart struct {
	PartId   uint `json:"partId"`
	SortCode uint `json:"sortCode"`
}

type SortPage struct {
	PageId   uint `json:"pageId"`
	SortCode uint `json:"sortCode"`
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
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SaveFail))
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
	pageId, pageIdErr := strconv.ParseUint(c.Param("pageId"), 10, 32)
	editable, editableErr := strconv.ParseBool(c.DefaultQuery("editable", "false"))
	if pageIdErr != nil || editableErr != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPage(uint(pageId), editable)))
	}
}

func savePageHandler(c *gin.Context) {
	pageDto := &PageDto{}
	if err := c.ShouldBind(pageDto); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if savePage(pageDto) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func editPageHandler(c *gin.Context) {
	pageDto := &PageDto{}
	if err := c.ShouldBind(pageDto); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if editPage(pageDto) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(UpdateFail))
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

func getHistoricalPagesHandler(c *gin.Context) {
	if pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getHistoricalPages(uint(pageId))))
	}
}

func getTagsHandler(c *gin.Context) {
	Result(c, OkData(getTags()))
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

func sortBooksHandler(c *gin.Context) {
	sortedBooks := &[]SortBook{}
	if err := c.ShouldBind(sortedBooks); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortBooks(sortedBooks); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}

func sortPartsHandler(c *gin.Context) {
	sortedParts := &[]SortPart{}
	if err := c.ShouldBind(sortedParts); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortParts(sortedParts); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}

func sortPagesHandler(c *gin.Context) {
	sortedPages := &[]SortPage{}
	if err := c.ShouldBind(sortedPages); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortPages(sortedPages); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}
