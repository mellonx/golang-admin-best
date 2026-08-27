package middleware

import (
	"golang-admin-best/pkg/response"
	"golang-admin-best/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			return
		}

		// 兼容两种格式：
		// 1. 直接 token（前端当前使用）
		// 2. "Bearer {token}" 标准格式
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				tokenString = parts[1]
			}
		}

		// 验证 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "令牌无效或已过期")
			return
		}

		// 存储用户信息到上下文
		c.Set("userId", claims.UserID)
		c.Set("userName", claims.UserName)

		c.Next()
	}
}
