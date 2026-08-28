package service

import (
	"golang-admin-best/internal/repository"
)

type SystemService struct {
	repo *repository.Repository
}

func NewSystemService(repo *repository.Repository) *SystemService {
	return &SystemService{repo: repo}
}

// PageResult 通用分页结果（对应前端 Api.Common.PaginatedResponse）
type PageResult struct {
	Records interface{} `json:"records"`
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
}

// UserListItem 用户列表项（对应前端 Api.SystemManage.UserListItem）
type UserListItem struct {
	ID         uint     `json:"id"`
	Avatar     string   `json:"avatar"`
	Status     string   `json:"status"`
	UserName   string   `json:"userName"`
	UserGender string   `json:"userGender"`
	NickName   string   `json:"nickName"`
	UserPhone  string   `json:"userPhone"`
	UserEmail  string   `json:"userEmail"`
	UserRoles  []string `json:"userRoles"`
	CreateBy   string   `json:"createBy"`
	CreateTime string   `json:"createTime"`
	UpdateBy   string   `json:"updateBy"`
	UpdateTime string   `json:"updateTime"`
}

// GetUserList 分页获取用户列表
func (s *SystemService) GetUserList(current, size int, userName string) (*PageResult, error) {
	offset := (current - 1) * size
	users, total, err := s.repo.User.FindPage(offset, size, userName)
	if err != nil {
		return nil, err
	}

	records := make([]UserListItem, 0, len(users))
	for _, u := range users {
		roles := make([]string, 0, len(u.Roles))
		for _, r := range u.Roles {
			roles = append(roles, r.RoleCode)
		}
		status := "1" // 1:启用
		if u.Status != 1 {
			status = "2" // 2:禁用
		}
		records = append(records, UserListItem{
			ID:         u.ID,
			Avatar:     u.Avatar,
			Status:     status,
			UserName:   u.UserName,
			NickName:   u.UserName,
			UserEmail:  u.Email,
			UserRoles:  roles,
			CreateTime: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &PageResult{Records: records, Current: current, Size: size, Total: total}, nil
}

// RoleListItem 角色列表项（对应前端 Api.SystemManage.RoleListItem）
type RoleListItem struct {
	RoleID      uint   `json:"roleId"`
	RoleName    string `json:"roleName"`
	RoleCode    string `json:"roleCode"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreateTime  string `json:"createTime"`
}

// GetRoleList 分页获取角色列表
func (s *SystemService) GetRoleList(current, size int, roleName string) (*PageResult, error) {
	offset := (current - 1) * size
	roles, total, err := s.repo.Role.FindPage(offset, size, roleName)
	if err != nil {
		return nil, err
	}

	records := make([]RoleListItem, 0, len(roles))
	for _, r := range roles {
		records = append(records, RoleListItem{
			RoleID:      r.ID,
			RoleName:    r.RoleName,
			RoleCode:    r.RoleCode,
			Description: r.Description,
			Enabled:     r.Status == 1,
			CreateTime:  r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &PageResult{Records: records, Current: current, Size: size, Total: total}, nil
}
