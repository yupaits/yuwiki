package yuwiki

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var sessionUserId uint

func newAuthMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := sessions.Default(c).Get(Config.SessionAuth)
		if userId != nil {
			userId := userId.(uint)
			if userId > 0 {
				sessionUserId = userId
				c.Next()
				return
			}
		}
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			//异步请求返回401状态码
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			//重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
		}
	}
}

func adminMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, err := getCurrentUser(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, MsgFail(err.Error()))
		} else if user.Admin {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, CodeFail(http.StatusForbidden))
		}
	}
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}
		log.SetFormatter(&logrus.TextFormatter{DisableColors: false, TimestampFormat: DateTimeLogLayout})
		log.WithFields(logrus.Fields{
			"status":   fmt.Sprintf("%d", statusCode),
			"latency":  fmt.Sprintf("%v", latency),
			"clientIP": fmt.Sprintf("%s", clientIP),
			"method":   fmt.Sprintf("%s", method),
			"path":     fmt.Sprintf("%s", path),
			"comment":  fmt.Sprintf("%s", comment),
		}).Info("GinLogger")
	}
}

func getUserId() uint {
	return sessionUserId
}
