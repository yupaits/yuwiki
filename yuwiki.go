package yuwiki

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

var Config *AppConfig

func Run() {
	Config = InitConfig()

	if Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if Config.LogFile != "" {
		logFile, _ := os.OpenFile(Config.LogFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	}

	r := gin.Default()

	r.Static("/static", Config.Http.StaticPath)
	r.StaticFile("/favicon.ico", Config.Http.Favicon)
	r.LoadHTMLGlob(Config.Http.HtmlPathPattern)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	books := r.Group("/books")
	{
		books.GET("", GetBooksHandler)
		books.GET("/{bookId}/parts", GetBookPartsHandler)
		books.GET("/shared", GetSharedBooksHandler)
		books.POST("", AddBookHandler)
		books.PUT("/{bookId}", EditBookHandler)
		books.DELETE("/{bookId}", DeleteBookHandler)
		books.POST("/share", ShareBookHandler)
	}

	parts := r.Group("/parts")
	{
		parts.GET("/{partId}/pages", GetPartPagesHandler)
		parts.POST("", AddPartHandler)
		parts.PUT("/{partId}", EditPartHandler)
		parts.DELETE("/{partId}", DeletePartHandler)
	}

	pages := r.Group("/pages")
	{
		pages.GET("/{pageId}", GetPageHandler)
		pages.POST("", AddPageHandler)
		pages.PUT("/{pageId}", EditPageHandler)
		pages.DELETE("/{pageId}", DeletePageHandler)
	}

	user := r.Group("/user")
	{
		user.GET("", GetUserInfoHandler)
		user.PUT("/modify-password", ModifyPasswordHandler)
	}

	r.GET("/star/items", GetStarItemsHandler)
	r.POST("/site/search", SiteSearchHandler)

	log.Fatalln(r.Run(":" + Config.Http.Port))
}
