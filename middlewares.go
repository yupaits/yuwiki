package yuwiki

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sessionUserId uint

func NewAuthMid() gin.HandlerFunc {
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

func getUserId() uint {
	return sessionUserId
}
