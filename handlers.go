package yuwiki

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetBooksHandler(c *gin.Context) {
	Result(c, OkData(GetBooks()))
}

func GetBookPartsHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(GetBookParts(uint(bookId))))
	}
}

func GetSharedBooksHandler(c *gin.Context) {

}

func SaveBookHandler(c *gin.Context) {
	book := Book{}
	if err := c.ShouldBindJSON(book); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if SaveBook(&book) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(CreateFail))
	}
}

func DeleteBookHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if DeleteBook(uint(bookId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func ShareBookHandler(c *gin.Context) {

}

func GetPartPagesHandler(c *gin.Context) {

}

func SavePartHandler(c *gin.Context) {
	part := Part{}
	if err := c.ShouldBind(part); err != nil {
		Result(c, CodeFail(ParamsError))
		return
	}
	if part.Protected && part.Password == "" {
		Result(c, CodeFail(ParamsError))
		return
	}
	if SavePart(&part) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(CreateFail))
	}
}

func DeletePartHandler(c *gin.Context) {

}

func GetPageHandler(c *gin.Context) {

}

func SavePageHandler(c *gin.Context) {

}

func DeletePageHandler(c *gin.Context) {

}

func GetUserInfoHandler(c *gin.Context) {

}

func ModifyPasswordHandler(c *gin.Context) {

}

func GetStarItemsHandler(c *gin.Context) {

}

func SiteSearchHandler(c *gin.Context) {

}
