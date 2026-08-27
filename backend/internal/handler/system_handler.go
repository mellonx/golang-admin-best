package handler

import (
	"art-design-pro-api/internal/service"
	"art-design-pro-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	menuService *service.MenuService
}

func NewSystemHandler(menuService *service.MenuService) *SystemHandler {
	return &SystemHandler{menuService: menuService}
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
