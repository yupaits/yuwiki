package yuwiki

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var sessionUser *User

func NewAuthMid(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c)
		userObj := sessions.Default(c).Get(Config.SessionAuth)
		if userObj != nil {
			user := userObj.(*User)
			if user.ID > 0 {
				sessionUser = user
				c.Next()
			}
		} else {
			//TODO
			//c.Redirect(http.StatusMovedPermanently, "/login")
			//c.Abort()
		}
	}
}

func getUserId() uint {
	return sessionUser.ID
}
