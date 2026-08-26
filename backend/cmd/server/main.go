package main

import (
	"art-design-pro-api/internal/middleware"
	"art-design-pro-api/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 引擎
	r := gin.Default()

	// 应用 CORS 中间件
	r.Use(middleware.CORS())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "healthy",
			"service": "art-design-pro-api",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				// TODO: 实现登录逻辑，见 docs/golang-code.md
				response.Error(c, 500, "登录接口待实现，请参考 docs/golang-code.md")
			})
		}

		// 用户相关
		user := api.Group("/user")
		{
			user.GET("/info", func(c *gin.Context) {
				// TODO: 实现获取用户信息，见 docs/golang-code.md
				response.Error(c, 500, "用户信息接口待实现，请参考 docs/golang-code.md")
			})
		}

		// 系统菜单
		v3 := api.Group("/v3")
		{
			system := v3.Group("/system")
			{
				system.GET("/menus", func(c *gin.Context) {
					// TODO: 实现菜单树接口，见 docs/golang-code.md
					response.Error(c, 500, "菜单接口待实现，请参考 docs/golang-code.md")
				})
			}
		}
	}

	// 启动服务
	log.Println("🚀 服务启动在 http://localhost:8080")
	log.Println("📖 完整实现请参考 docs/golang-code.md")
	log.Println("🩺 健康检查: http://localhost:8080/health")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("启动服务失败:", err)
	}
}
