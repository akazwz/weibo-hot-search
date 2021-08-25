package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"` //omitempty nil or default
	Msg  string      `json:"msg,omitempty"`
}

const (
	SUCCESS  = 2000
	ERROR    = 4000
	PROGRESS = 2020
)

// NotFound 路由不存在
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code: ERROR,
		Msg:  "404 not found",
	})
}

// Teapot 418 teapot
func Teapot(c *gin.Context) {
	c.JSON(http.StatusTeapot, gin.H{
		"message": "I'm a teapot",
		"story": "This code was defined in 1998 " +
			"as one of the traditional IETF April Fools' jokes," +
			" in RFC 2324, Hyper Text Coffee Pot Control Protocol," +
			" and is not expected to be implemented by actual HTTP servers." +
			" However, known implementations do exist.",
	})
}

func CommonFailed(code int, message string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code: code,
		Msg:  message,
	})
}

func CommonSuccess(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

func Created(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusCreated, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

func DeletedNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
