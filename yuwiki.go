package yuwiki

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte(Config.Secret))
	store.Options(sessions.Options{
		MaxAge:   86400,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	sessionMid := sessions.Sessions(Config.SessionCookie, store)
	authorize := NewAuthMid(sessionMid)

	r.Static("/static", Config.Http.StaticPath)
	r.StaticFile("/favicon.ico", Config.Http.Favicon)
	//r.LoadHTMLGlob(Config.Http.HtmlPathPattern)
	r.LoadHTMLFiles("./templates/login.html", "./ui/dist/index.html")

	r.GET("/logout", authorize, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete(Config.SessionAuth)
		if err := session.Save(); err != nil {

		}
	})

	r.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		if ok, user, err := checkLogin(c); ok {
			session.Set(Config.SessionAuth, user)
			if err := session.Save(); err != nil {
				//TODO
			}
		} else if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		}

	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	books := r.Group("/books").Use(authorize)
	{
		books.GET("", getBooksHandler)
		books.GET("/:bookId/parts", getBookPartsHandler)
		books.POST("", saveBookHandler)
		books.PUT("/:bookId", saveBookHandler)
		books.DELETE("/:bookId", deleteBookHandler)
		books.POST("/share", shareBookHandler)
		books.POST("/sort", sortBooksHandler)
	}

	parts := r.Group("/parts").Use(authorize)
	{
		parts.GET("/:partId/pages", getPartPagesHandler)
		parts.GET("/:partId", getPartHandler)
		parts.POST("", savePartHandler)
		parts.PUT("/:partId", savePartHandler)
		parts.DELETE("/:partId", deletePartHandler)
		parts.POST("/sort", sortPartsHandler)
	}

	pages := r.Group("/pages").Use(authorize)
	{
		pages.GET("/:pageId", getPageHandler)
		pages.GET("/:pageId/history", getHistoricalPagesHandler)
		pages.POST("", savePageHandler)
		pages.PUT("/:pageId", savePageHandler)
		pages.PUT("/:pageId/edit", editPageHandler)
		pages.DELETE("/:pageId", deletePageHandler)
		pages.POST("/sort", sortPagesHandler)
	}

	tags := r.Group("/tags").Use(authorize)
	{
		tags.GET("", getTagsHandler)
	}

	user := r.Group("/user").Use(authorize)
	{
		user.GET("", getUserInfoHandler)
		user.PUT("/modify-password", modifyPasswordHandler)
	}

	r.GET("/shared/books", authorize, getSharedBooksHandler)
	r.GET("/star/items", authorize, getStarItemsHandler)
	r.POST("/site/search", authorize, siteSearchHandler)

	log.Fatal(r.Run(":" + Config.Http.Port))
}
