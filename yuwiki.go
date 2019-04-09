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

	existsDir := Mkdirs(Config.LogFile)
	logFile, err := os.OpenFile(Config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	hasLogFile := existsDir && err == nil
	if hasLogFile {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	}

	InitDb(Config.DataSource.DdlUpdate)

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
		books.GET("", getBooksHandler)
		books.GET("/:bookId/parts", getBookPartsHandler)
		books.POST("", saveBookHandler)
		books.PUT("/:bookId", saveBookHandler)
		books.DELETE("/:bookId", deleteBookHandler)
		books.POST("/share", shareBookHandler)
	}

	parts := r.Group("/parts")
	{
		parts.GET("/:partId/pages", getPartPagesHandler)
		parts.GET("/:partId", getPartHandler)
		parts.POST("", savePartHandler)
		parts.PUT("/:partId", savePartHandler)
		parts.DELETE("/:partId", deletePartHandler)
	}

	pages := r.Group("/pages")
	{
		pages.GET("/:pageId", getPageHandler)
		pages.GET("/:pageId/history", getHistoricalPagesHandler)
		pages.POST("", savePageHandler)
		pages.PUT("/:pageId", savePageHandler)
		pages.PUT("/:pageId/edit", editPageHandler)
		pages.DELETE("/:pageId", deletePageHandler)
	}

	user := r.Group("/user")
	{
		user.GET("", getUserInfoHandler)
		user.PUT("/modify-password", modifyPasswordHandler)
	}

	r.GET("/shared/books", getSharedBooksHandler)
	r.GET("/star/items", getStarItemsHandler)
	r.POST("/site/search", siteSearchHandler)

	log.Fatal(r.Run(":" + Config.Http.Port))
}
