package router

import (
	"golang-admin-best/internal/handler"
	"golang-admin-best/internal/middleware"
	"golang-admin-best/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handlers 聚合所有 HTTP 处理器，便于统一注入路由层
type Handlers struct {
	Auth   *handler.AuthHandler
	System *handler.SystemHandler
}

// Setup 创建 Gin 引擎并注册所有路由
func Setup(h *Handlers) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", healthCheck)

	// 注册 API 路由
	registerAPIRoutes(r, h)

	return r
}

// healthCheck 健康检查处理器
func healthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "healthy",
		"service": "golang-admin-best",
	})
}

// registerAPIRoutes 注册 /api 下的所有路由
func registerAPIRoutes(r *gin.Engine, h *Handlers) {
	api := r.Group("/api")

	// 公开接口（无需认证）
	registerPublicRoutes(api, h)

	// 需要认证的接口
	registerAuthRoutes(api, h)
}

// registerPublicRoutes 注册公开路由
func registerPublicRoutes(api *gin.RouterGroup, h *Handlers) {
	api.POST("/auth/login", h.Auth.Login)
}

// registerAuthRoutes 注册需要 JWT 认证的路由
func registerAuthRoutes(api *gin.RouterGroup, h *Handlers) {
	auth := api.Group("")
	auth.Use(middleware.JWTAuth())

	// 用户相关
	auth.GET("/user/info", h.Auth.GetUserInfo)
	auth.GET("/user/list", h.System.GetUserList)

	// 角色相关
	auth.GET("/role/list", h.System.GetRoleList)

	// 系统管理
	v3 := auth.Group("/v3")
	v3.GET("/system/menus", h.System.GetMenuList)
}
