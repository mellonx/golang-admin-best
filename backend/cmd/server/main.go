package main

import (
	"art-design-pro-api/internal/config"
	"art-design-pro-api/internal/database"
	"art-design-pro-api/internal/handler"
	"art-design-pro-api/internal/middleware"
	"art-design-pro-api/internal/repository"
	"art-design-pro-api/internal/service"
	"art-design-pro-api/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置（优先 .env，否则环境变量）
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}
	log.Printf("⚙️  Config loaded (mode: %s, port: %s)\n",
		config.AppConfig.Server.Mode, config.AppConfig.Server.Port)

	// 2. 初始化数据库连接
	if err := database.Init(); err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer database.Close()

	// 3. 初始化 Repository 层（数据访问抽象）
	repo := repository.NewRepository(database.Get()) // 使用默认连接

	// 4. 初始化 Service 层（业务逻辑）
	authService := service.NewAuthService(repo)
	menuService := service.NewMenuService(repo)

	// 5. 初始化 Handler 层（HTTP 处理）
	authHandler := handler.NewAuthHandler(authService)
	systemHandler := handler.NewSystemHandler(menuService)

	// 6. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 7. 创建路由
	r := gin.Default()
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "healthy",
			"service": "art-design-pro-api",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 公开接口（无需认证）
		api.POST("/auth/login", authHandler.Login)

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/info", authHandler.GetUserInfo)
			auth.GET("/v3/system/menus", systemHandler.GetMenuList)
		}
	}

	// 8. 启动服务
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s\n", addr)
	log.Printf("🩺 Health check: http://localhost%s/health\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
