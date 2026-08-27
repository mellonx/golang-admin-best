package main

import (
	"art-design-pro-api/internal/config"
	"art-design-pro-api/internal/database"
	"art-design-pro-api/internal/middleware"
	"art-design-pro-api/internal/repository"
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
	repo := repository.NewRepository(database.DB)
	_ = repo // TODO: 注入到 service/handler 层

	// 4. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 5. 创建路由
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
		// TODO: 使用 repo 实现真实业务逻辑，见 docs/golang-code.md
		api.POST("/auth/login", func(c *gin.Context) {
			response.Error(c, 500, "登录接口待实现，请参考 docs/golang-code.md")
		})

		api.GET("/user/info", func(c *gin.Context) {
			response.Error(c, 500, "用户信息接口待实现，请参考 docs/golang-code.md")
		})

		v3 := api.Group("/v3")
		{
			v3.GET("/system/menus", func(c *gin.Context) {
				response.Error(c, 500, "菜单接口待实现，请参考 docs/golang-code.md")
			})
		}
	}

	// 6. 启动服务
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s\n", addr)
	log.Printf("🩺 Health check: http://localhost%s/health\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
