package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式，对应前端 BaseResponse
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 状态码常量，对应前端 ApiStatus
const (
	CodeSuccess      = 200 // 成功
	CodeError        = 400 // 请求错误
	CodeUnauthorized = 401 // 未授权（前端自动登出）
	CodeForbidden    = 403 // 无权限
	CodeServerError  = 500 // 服务器错误
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 带自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应（HTTP 200，业务码在 body）
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// Unauthorized 401 未授权，前端会自动登出
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: CodeUnauthorized,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
}
