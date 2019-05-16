package yuwiki

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	FAIL              = 10000
	LoginFail         = 10001
	InvalidSession    = 10002
	ParamsError       = 11001
	DataNotFound      = 11002
	DataCannotDelete  = 11003
	DataConflict      = 11004
	CreateFail        = 20001
	UpdateFail        = 20002
	DeleteFail        = 20003
	SaveFail          = 20004
	SortFail          = 20005
	FileUploadFail    = 30001
	CaptchaVerifyFail = 30002
)

var messages = map[int]string{
	FAIL:              "系统异常",
	LoginFail:         "用户名或密码错误",
	InvalidSession:    "无效的Session信息",
	ParamsError:       "参数校验失败",
	DataNotFound:      "查找的数据不存在或已被删除",
	DataCannotDelete:  "数据无法删除",
	DataConflict:      "不允许的数据重复",
	CreateFail:        "创建记录失败",
	UpdateFail:        "更新记录失败",
	DeleteFail:        "删除记录失败",
	SaveFail:          "保存记录失败",
	SortFail:          "排序操作失败",
	FileUploadFail:    "文件上传失败",
	CaptchaVerifyFail: "验证码校验失败",
}

func Ok() gin.H {
	return gin.H{"code": http.StatusOK, "msg": http.StatusText(http.StatusOK), "success": true}
}

func OkData(data interface{}) gin.H {
	return gin.H{"code": http.StatusOK, "msg": http.StatusText(http.StatusOK), "success": true, "data": data}
}

func Fail() gin.H {
	return gin.H{"code": FAIL, "msg": messages[FAIL], "success": false}
}

func CodeFail(code int) gin.H {
	if msg, ok := messages[code]; ok {
		return gin.H{"code": code, "msg": msg, "success": false}
	} else if http.StatusText(code) != "" {
		return gin.H{"code": code, "msg": http.StatusText(code), "success": false}
	}
	return gin.H{"code": FAIL, "msg": messages[FAIL], "success": false}
}

func FailData(data interface{}) gin.H {
	return gin.H{"code": FAIL, "msg": messages[FAIL], "success": false, "data": data}
}

func CodeFailData(code int, data interface{}) gin.H {
	if msg, ok := messages[code]; ok {
		return gin.H{"code": code, "msg": msg, "success": false, "data": data}
	} else if http.StatusText(code) != "" {
		return gin.H{"code": code, "msg": http.StatusText(code), "success": false, "data": data}
	}
	return gin.H{"code": FAIL, "msg": messages[FAIL], "success": false, "data": data}
}

func MsgFail(msg string) gin.H {
	return gin.H{"code": FAIL, "msg": msg, "success": false}
}

func MsgFailData(msg string, data interface{}) gin.H {
	return gin.H{"code": FAIL, "msg": msg, "success": false, "data": data}
}

func Result(c *gin.Context, res gin.H) {
	c.JSON(http.StatusOK, res)
}
