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
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func getUserId() uint {
	return sessionUserId
}
