package yuwiki

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PasswordModify struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"oldPassword" binding:"required"`
	ConfirmPassword string `json:"oldPassword" binding:"required"`
}

func getBooksHandler(c *gin.Context) {
	Result(c, OkData(getBooks()))
}

func GetBookPartsHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getBookParts(uint(bookId))))
	}
}

func getSharedBooksHandler(c *gin.Context) {

}

func saveBookHandler(c *gin.Context) {
	book := Book{}
	if err := c.ShouldBindJSON(book); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if saveBook(&book) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(CreateFail))
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

}

func getPartPagesHandler(c *gin.Context) {

}

func savePartHandler(c *gin.Context) {
	part := Part{}
	if err := c.ShouldBind(part); err != nil {
		Result(c, CodeFail(ParamsError))
		return
	}
	if part.Protected && part.Password == "" {
		Result(c, CodeFail(ParamsError))
		return
	}
	if ok, err := savePart(&part); ok {
		Result(c, Ok())
	} else {
		Result(c, MsgFail(err.Error()))
	}
}

func deletePartHandler(c *gin.Context) {

}

func getPageHandler(c *gin.Context) {

}

func savePageHandler(c *gin.Context) {

}

func deletePageHandler(c *gin.Context) {

}

func getUserInfoHandler(c *gin.Context) {

}

func modifyPasswordHandler(c *gin.Context) {

}

func getStarItemsHandler(c *gin.Context) {

}

func siteSearchHandler(c *gin.Context) {

}
