package handler

import (
	"golang-admin-best/internal/service"
	"golang-admin-best/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginParams 登录参数（对应前端 Api.Auth.LoginParams）
type LoginParams struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var params LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, response.CodeError, "参数错误")
		return
	}

	result, err := h.authService.Login(params.UserName, params.Password)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "登录成功", result)
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userId")

	info, err := h.authService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, info)
}
