# Art Design Pro - Golang 后端核心代码

基于 **Gin + GORM + JWT** 实现。

## 一、依赖安装

```bash
go mod init golang-admin-best

go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

## 二、统一响应结构

```go
// utils/response.go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构（对应前端 BaseResponse）
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 前端 ApiStatus 状态码
const (
	CodeSuccess      = 200 // 成功
	CodeUnauthorized = 401 // 未授权（触发前端登出）
	CodeForbidden    = 403 // 无权限
	CodeError        = 500 // 服务器错误
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Fail(c *gin.Context, code int, msg string) {
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
```

## 三、JWT 工具

```go
// utils/jwt.go
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// GenerateToken 生成访问令牌
func GenerateToken(userID uint, userName string) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, userName string) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
```

## 四、数据模型

```go
// models/models.go
package models

import "time"

// User 用户
type User struct {
	ID        uint      `gorm:"primaryKey" json:"userId"`
	UserName  string    `gorm:"uniqueIndex;size:50" json:"userName"`
	Password  string    `gorm:"size:255" json:"-"` // 不返回给前端
	Email     string    `gorm:"size:100" json:"email"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int8      `gorm:"default:1" json:"status"`
	Roles     []Role    `gorm:"many2many:user_roles;" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Role 角色
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"roleId"`
	RoleCode    string       `gorm:"uniqueIndex;size:50" json:"roleCode"`
	RoleName    string       `gorm:"size:50" json:"roleName"`
	Description string       `gorm:"size:255" json:"description"`
	Status      int8         `gorm:"default:1" json:"status"`
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"-"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
	CreatedAt   time.Time    `json:"createTime"`
}

// Menu 菜单
type Menu struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ParentID  uint         `gorm:"default:0;index" json:"-"`
	Path      string       `gorm:"size:255" json:"path"`
	Name      string       `gorm:"size:100" json:"name"`
	Component string       `gorm:"size:255" json:"component"`
	Redirect  string       `gorm:"size:255" json:"redirect,omitempty"`
	Title     string       `gorm:"size:100" json:"-"`
	Icon      string       `gorm:"size:100" json:"-"`
	KeepAlive int8         `gorm:"default:0" json:"-"`
	IsHide    int8         `gorm:"default:0" json:"-"`
	IsIframe  int8         `gorm:"default:0" json:"-"`
	Link      string       `gorm:"size:255" json:"-"`
	Sort      int          `gorm:"default:0" json:"-"`
	Status    int8         `gorm:"default:1" json:"-"`
	Permissions []Permission `gorm:"foreignKey:MenuID" json:"-"`
}

// Permission 按钮级权限
type Permission struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	MenuID   uint   `gorm:"index" json:"-"`
	Title    string `gorm:"size:50" json:"title"`
	AuthMark string `gorm:"size:50" json:"authMark"`
}
```

## 五、菜单 DTO（转换为前端需要的格式）

```go
// models/dto.go
package models

// MenuMeta 菜单元数据（对应前端 RouteMeta）
type MenuMeta struct {
	Title     string     `json:"title"`
	Icon      string     `json:"icon,omitempty"`
	KeepAlive bool       `json:"keepAlive,omitempty"`
	IsHide    bool       `json:"isHide,omitempty"`
	IsIframe  bool       `json:"isIframe,omitempty"`
	Link      string     `json:"link,omitempty"`
	Roles     []string   `json:"roles,omitempty"`
	AuthList  []AuthItem `json:"authList,omitempty"`
}

// AuthItem 操作权限项
type AuthItem struct {
	Title    string `json:"title"`
	AuthMark string `json:"authMark"`
}

// MenuTree 菜单树（对应前端 AppRouteRecord）
type MenuTree struct {
	ID        uint       `json:"id"`
	Path      string     `json:"path"`
	Name      string     `json:"name"`
	Component string     `json:"component,omitempty"`
	Redirect  string     `json:"redirect,omitempty"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuTree `json:"children,omitempty"`
}
```

## 六、认证控制器

```go
// controllers/auth_controller.go
package controllers

import (
	"golang-admin-best/models"
	"golang-admin-best/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// LoginParams 登录参数（对应前端 Api.Auth.LoginParams）
type LoginParams struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (ac *AuthController) Login(c *gin.Context) {
	var params LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.Fail(c, utils.CodeError, "参数错误")
		return
	}

	// 查询用户
	var user models.User
	if err := ac.DB.Where("user_name = ?", params.UserName).First(&user).Error; err != nil {
		utils.Fail(c, utils.CodeError, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		utils.Fail(c, utils.CodeError, "用户名或密码错误")
		return
	}

	// 检查状态
	if user.Status != 1 {
		utils.Fail(c, utils.CodeError, "账号已被禁用")
		return
	}

	// 生成 token
	token, _ := utils.GenerateToken(user.ID, user.UserName)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID, user.UserName)

	// 返回（对应前端 Api.Auth.LoginResponse）
	utils.SuccessWithMsg(c, "登录成功", gin.H{
		"token":        token,
		"refreshToken": refreshToken,
	})
}

// GetUserInfo 获取用户信息
func (ac *AuthController) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userId") // 从中间件获取

	var user models.User
	if err := ac.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		utils.Fail(c, utils.CodeError, "获取用户信息失败")
		return
	}

	// 提取角色代码
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.RoleCode)
	}

	// 提取用户拥有的所有按钮权限
	buttons := ac.getUserButtons(user.ID)

	// 返回（对应前端 Api.Auth.UserInfo）
	utils.Success(c, gin.H{
		"userId":   user.ID,
		"userName": user.UserName,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"roles":    roles,
		"buttons":  buttons,
	})
}

// getUserButtons 获取用户所有按钮权限
func (ac *AuthController) getUserButtons(userID uint) []string {
	var authMarks []string
	ac.DB.Table("permissions").
		Select("DISTINCT permissions.auth_mark").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.auth_mark", &authMarks)

	if authMarks == nil {
		authMarks = []string{}
	}
	return authMarks
}
```

## 七、菜单服务（核心：构建权限菜单树）

```go
// services/menu_service.go
package services

import (
	"golang-admin-best/models"

	"gorm.io/gorm"
)

type MenuService struct {
	DB *gorm.DB
}

// GetUserMenus 获取用户有权限的菜单树
func (ms *MenuService) GetUserMenus(userID uint) []models.MenuTree {
	// 1. 查询用户的所有角色ID
	var roleIDs []uint
	ms.DB.Table("user_roles").Where("user_id = ?", userID).Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return []models.MenuTree{}
	}

	// 2. 查询这些角色关联的菜单ID（去重）
	var menuIDs []uint
	ms.DB.Table("role_menus").
		Where("role_id IN ?", roleIDs).
		Distinct("menu_id").
		Pluck("menu_id", &menuIDs)

	if len(menuIDs) == 0 {
		return []models.MenuTree{}
	}

	// 3. 查询菜单详情
	var menus []models.Menu
	ms.DB.Where("id IN ? AND status = 1", menuIDs).
		Order("sort ASC").
		Find(&menus)

	// 4. 查询这些菜单的按钮权限（只返回该角色拥有的权限）
	permMap := ms.getMenuPermissions(roleIDs, menuIDs)

	// 5. 构建菜单树
	return ms.buildMenuTree(menus, 0, permMap)
}

// getMenuPermissions 获取菜单对应的按钮权限（按角色过滤）
func (ms *MenuService) getMenuPermissions(roleIDs []uint, menuIDs []uint) map[uint][]models.AuthItem {
	type permResult struct {
		MenuID   uint
		Title    string
		AuthMark string
	}
	var results []permResult

	ms.DB.Table("permissions").
		Select("permissions.menu_id, permissions.title, permissions.auth_mark").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ? AND permissions.menu_id IN ?", roleIDs, menuIDs).
		Distinct("permissions.id, permissions.menu_id, permissions.title, permissions.auth_mark").
		Scan(&results)

	permMap := make(map[uint][]models.AuthItem)
	for _, r := range results {
		permMap[r.MenuID] = append(permMap[r.MenuID], models.AuthItem{
			Title:    r.Title,
			AuthMark: r.AuthMark,
		})
	}
	return permMap
}

// buildMenuTree 递归构建菜单树
func (ms *MenuService) buildMenuTree(menus []models.Menu, parentID uint, permMap map[uint][]models.AuthItem) []models.MenuTree {
	tree := make([]models.MenuTree, 0)

	for _, menu := range menus {
		if menu.ParentID != parentID {
			continue
		}

		node := models.MenuTree{
			ID:        menu.ID,
			Path:      menu.Path,
			Name:      menu.Name,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: models.MenuMeta{
				Title:     menu.Title,
				Icon:      menu.Icon,
				KeepAlive: menu.KeepAlive == 1,
				IsHide:    menu.IsHide == 1,
				IsIframe:  menu.IsIframe == 1,
				Link:      menu.Link,
				AuthList:  permMap[menu.ID],
			},
		}

		// 递归构建子菜单
		children := ms.buildMenuTree(menus, menu.ID, permMap)
		if len(children) > 0 {
			node.Children = children
		}

		tree = append(tree, node)
	}

	return tree
}
```

## 八、系统控制器

```go
// controllers/system_controller.go
package controllers

import (
	"golang-admin-best/services"
	"golang-admin-best/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemController struct {
	DB          *gorm.DB
	MenuService *services.MenuService
}

// GetMenuList 获取菜单列表（后端模式）
func (sc *SystemController) GetMenuList(c *gin.Context) {
	userID := c.GetUint("userId")
	menus := sc.MenuService.GetUserMenus(userID)
	utils.Success(c, menus)
}
```

## 九、认证中间件

```go
// middleware/auth.go
package middleware

import (
	"golang-admin-best/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证格式错误")
			return
		}

		// 验证 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "令牌无效或已过期")
			return
		}

		// 存储用户信息到上下文
		c.Set("userId", claims.UserID)
		c.Set("userName", claims.UserName)

		c.Next()
	}
}
```

## 十、CORS 中间件

```go
// middleware/cors.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
```

## 十一、路由注册

```go
// routes/routes.go
package routes

import (
	"golang-admin-best/controllers"
	"golang-admin-best/middleware"
	"golang-admin-best/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	// 初始化服务和控制器
	menuService := &services.MenuService{DB: db}
	authController := &controllers.AuthController{DB: db}
	systemController := &controllers.SystemController{
		DB:          db,
		MenuService: menuService,
	}

	api := r.Group("/api")
	{
		// 公开接口（无需认证）
		api.POST("/auth/login", authController.Login)

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/info", authController.GetUserInfo)
			auth.GET("/v3/system/menus", systemController.GetMenuList)
		}
	}

	return r
}
```

## 十二、主程序

```go
// main.go
package main

import (
	"golang-admin-best/models"
	"golang-admin-best/routes"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:password@tcp(127.0.0.1:3306)/golang_admin_best?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Permission{},
	)

	// 启动服务
	r := routes.SetupRoutes(db)
	log.Println("服务启动在 :8080")
	r.Run(":8080")
}
```

## 十三、前端对接配置

修改 `art-design-pro/.env.development`：

```bash
# 切换到后端模式
# 修改 .env 中的 VITE_ACCESS_MODE = backend

# API 代理指向你的 Golang 服务
VITE_API_PROXY_URL = http://localhost:8080
```

修改 `art-design-pro/.env`：
```bash
VITE_ACCESS_MODE = backend
```

## 十四、关键对接要点

### 1. Token 传递
前端在请求头携带：`Authorization: Bearer {token}`
（需检查前端 `src/utils/http` 的拦截器实现）

### 2. 响应格式统一
所有接口返回：`{ code, msg, data }`
- `code: 200` 成功
- `code: 401` 未授权（前端自动登出）

### 3. 菜单 component 路径
- 对应 `src/views` 下的文件路径
- 例：`/system/user` → `src/views/system/user/index.vue`
- 目录型菜单用 `/index/index` 或空字符串

### 4. 角色代码约定
- `R_SUPER`: 超级管理员
- `R_ADMIN`: 管理员
- `R_USER`: 普通用户
- 前端 `v-roles` 指令直接使用这些代码

### 5. 权限标识约定
- `add`, `edit`, `delete`, `export` 等
- 前端 `v-auth` 指令使用这些标识
