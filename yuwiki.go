package yuwiki

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Config *AppConfig

func Run() {
	Init()
	defer Db.Close()

	if Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger(), gin.Recovery())

	store := cookie.NewStore([]byte(Config.Secret))
	store.Options(sessions.Options{
		MaxAge: 86400,
		Path:   "/",
	})
	r.Use(sessions.Sessions(Config.SessionCookie, store))

	gob.Register(User{})
	authorize := newAuthMid()

	r.Static("/static", Config.Http.StaticPath)
	r.StaticFile("/favicon.ico", Config.Http.Favicon)
	r.LoadHTMLGlob(Config.Http.HtmlPathPattern)

	r.GET("/logout", authorize, func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(Config.SessionAuth)
		if userId == nil {
			Result(c, CodeFail(InvalidSession))
		} else {
			session.Clear()
			if err := session.Save(); err != nil {
				Result(c, MsgFail(err.Error()))
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
		}
	})

	r.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		if ok, user, err := checkLogin(c); ok && user != nil {
			session.Set(Config.SessionAuth, user.ID)
			if err := session.Save(); err != nil {
				Result(c, MsgFail(err.Error()))
			} else {
				Result(c, Ok())
			}
		} else if err != nil {
			Result(c, MsgFail(err.Error()))
		} else {
			Result(c, CodeFail(LoginFail))
		}
	})

	r.POST("/signup", signUpHandler)

	r.GET("/login", func(c *gin.Context) {
		userId := sessions.Default(c).Get(Config.SessionAuth)
		if userId != nil {
			userId := userId.(uint)
			if userId > 0 {
				c.Redirect(http.StatusFound, "/")
				return
			}
		}
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})
	r.GET("/", authorize, func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/index", authorize, func(c *gin.Context) {
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
		books.PUT("/:bookId/star", toggleStarBookHandler)
	}

	parts := r.Group("/parts").Use(authorize)
	{
		parts.GET("/:partId/pages", getPartPagesHandler)
		parts.GET("/:partId", getPartHandler)
		parts.POST("", savePartHandler)
		parts.PUT("/:partId", savePartHandler)
		parts.DELETE("/:partId", deletePartHandler)
		parts.POST("/sort", sortPartsHandler)
		parts.PUT("/:partId/star", toggleStarPartHandler)
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
		pages.PUT("/:pageId/star", toggleStarPageHandler)
	}

	tags := r.Group("/tags").Use(authorize)
	{
		tags.GET("", getTagsHandler)
	}

	user := r.Group("/user").Use(authorize)
	{
		user.GET("", getUserInfoHandler)
		user.PUT("/edit", editUserHandler)
		user.PUT("/modify-password", modifyPasswordHandler)
	}

	r.GET("/shared/books", authorize, getSharedBooksHandler)
	r.GET("/site/search", authorize, siteSearchHandler)

	r.POST("/upload", authorize, uploadFileHandler)
	r.GET("/file/:filename", authorize, downloadFileHandler)

	StartScheduler()

	log.WithField("port", Config.Http.Port).Info("yuwiki启动成功")
	log.Fatal(r.Run(":" + Config.Http.Port))
}

func Init() {
	//初始化配置信息
	Config = initConfig()
	//创建目录
	Mkdirs(Config.Path.BackupPath)
	Mkdirs(Config.Path.UploadPath)
	//初始化数据库
	InitDb()
	//配置日志
	InitLog()
}
