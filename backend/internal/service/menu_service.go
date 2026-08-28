package service

import (
	"golang-admin-best/internal/model"
	"golang-admin-best/internal/repository"
)

type MenuService struct {
	repo *repository.Repository
}

func NewMenuService(repo *repository.Repository) *MenuService {
	return &MenuService{repo: repo}
}

// GetUserMenus 获取用户有权限的菜单树
func (s *MenuService) GetUserMenus(userID uint) ([]model.MenuTree, error) {
	// 1. 查询用户的角色ID
	roles, err := s.repo.Role.FindByUserID(userID)
	if err != nil || len(roles) == 0 {
		return []model.MenuTree{}, nil
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}

	// 2. 查询角色关联的菜单（已按 sort 排序）
	menus, err := s.repo.Menu.FindByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	// 3. 查询菜单的按钮权限（按角色过滤）
	menuIDs := make([]uint, 0, len(menus))
	for _, m := range menus {
		menuIDs = append(menuIDs, m.ID)
	}

	permMap := make(map[uint][]model.AuthItem)
	if len(menuIDs) > 0 {
		perms, err := s.repo.Permission.FindByRoleIDsAndMenuIDs(roleIDs, menuIDs)
		if err != nil {
			return nil, err
		}
		for _, p := range perms {
			permMap[p.MenuID] = append(permMap[p.MenuID], model.AuthItem{
				Title:    p.Title,
				AuthMark: p.AuthMark,
			})
		}
	}

	// 4. 构建菜单树
	return s.buildMenuTree(menus, 0, permMap), nil
}

// buildMenuTree 递归构建菜单树
func (s *MenuService) buildMenuTree(menus []*model.Menu, parentID uint, permMap map[uint][]model.AuthItem) []model.MenuTree {
	tree := make([]model.MenuTree, 0)

	for _, menu := range menus {
		if menu.ParentID != parentID {
			continue
		}

		node := model.MenuTree{
			ID:        menu.ID,
			Path:      menu.Path,
			Name:      menu.Name,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: model.MenuMeta{
				Title:         menu.Title,
				Icon:          menu.Icon,
				KeepAlive:     menu.KeepAlive == 1,
				IsHide:        menu.IsHide == 1,
				IsHideTab:     menu.IsHideTab == 1,
				IsIframe:      menu.IsIframe == 1,
				IsFullPage:    menu.IsFullPage == 1,
				FixedTab:      menu.FixedTab == 1,
				Link:          menu.Link,
				ActivePath:    menu.ActivePath,
				ShowTextBadge: menu.ShowTextBadge,
				AuthList:      permMap[menu.ID],
			},
		}

		// 递归构建子菜单
		children := s.buildMenuTree(menus, menu.ID, permMap)
		if len(children) > 0 {
			node.Children = children
		}

		tree = append(tree, node)
	}

	return tree
}
