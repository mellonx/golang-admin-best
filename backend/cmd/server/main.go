package main

import (
	"golang-admin-best/internal/config"
	"golang-admin-best/internal/database"
	"golang-admin-best/internal/handler"
	"golang-admin-best/internal/repository"
	"golang-admin-best/internal/router"
	"golang-admin-best/internal/service"
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
	systemService := service.NewSystemService(repo)

	// 5. 初始化 Handler 层（HTTP 处理）
	handlers := &router.Handlers{
		Auth:   handler.NewAuthHandler(authService),
		System: handler.NewSystemHandler(menuService, systemService),
	}

	// 6. 设置 Gin 模式并注册路由
	gin.SetMode(config.AppConfig.Server.Mode)
	r := router.Setup(handlers)

	// 7. 启动服务
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s\n", addr)
	log.Printf("🩺 Health check: http://localhost%s/health\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
