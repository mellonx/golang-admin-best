package handler

import (
	"strconv"

	"golang-admin-best/internal/service"
	"golang-admin-best/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	menuService   *service.MenuService
	systemService *service.SystemService
}

func NewSystemHandler(menuService *service.MenuService, systemService *service.SystemService) *SystemHandler {
	return &SystemHandler{menuService: menuService, systemService: systemService}
}

// GetMenuList 获取当前用户的菜单树（后端权限模式核心接口）
func (h *SystemHandler) GetMenuList(c *gin.Context) {
	userID := c.GetUint("userId")

	menus, err := h.menuService.GetUserMenus(userID)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取菜单失败")
		return
	}

	response.Success(c, menus)
}

// GetUserList 分页获取用户列表
func (h *SystemHandler) GetUserList(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	size := parseIntQuery(c, "size", 20)
	userName := c.Query("userName")

	result, err := h.systemService.GetUserList(current, size, userName)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取用户列表失败")
		return
	}
	response.Success(c, result)
}

// GetRoleList 分页获取角色列表
func (h *SystemHandler) GetRoleList(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	size := parseIntQuery(c, "size", 20)
	roleName := c.Query("roleName")

	result, err := h.systemService.GetRoleList(current, size, roleName)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取角色列表失败")
		return
	}
	response.Success(c, result)
}

// parseIntQuery 解析整型查询参数，失败返回默认值
func parseIntQuery(c *gin.Context, key string, def int) int {
	if v := c.Query(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}
