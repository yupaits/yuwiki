package yuwiki

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBooksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, OkData([]Book{}))
}

func GetBookPartsHandler(c *gin.Context) {

}

func GetSharedBooksHandler(c *gin.Context) {

}

func AddBookHandler(c *gin.Context) {

}

func EditBookHandler(c *gin.Context) {

}

func DeleteBookHandler(c *gin.Context) {

}

func ShareBookHandler(c *gin.Context) {

}

func GetPartPagesHandler(c *gin.Context) {

}

func AddPartHandler(c *gin.Context) {

}

func EditPartHandler(c *gin.Context) {

}

func DeletePartHandler(c *gin.Context) {

}

func GetPageHandler(c *gin.Context) {

}

func AddPageHandler(c *gin.Context) {

}

func EditPageHandler(c *gin.Context) {

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
