package service

import (
	"art-design-pro-api/internal/repository"
	"art-design-pro-api/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 业务错误定义
var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("账号已被禁用")
	ErrUserNotFound       = errors.New("获取用户信息失败")
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

// LoginResult 登录结果
type LoginResult struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// Login 登录：校验用户名密码，生成令牌
func (s *AuthService) Login(userName, password string) (*LoginResult, error) {
	user, err := s.repo.User.FindByUserName(userName)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 校验状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 生成令牌
	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, RefreshToken: refreshToken}, nil
}

// UserInfoResult 用户信息结果（对应前端 Api.Auth.UserInfo）
type UserInfoResult struct {
	UserID   uint     `json:"userId"`
	UserName string   `json:"userName"`
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	Buttons  []string `json:"buttons"`
}

// GetUserInfo 获取用户信息：角色列表 + 按钮权限
func (s *AuthService) GetUserInfo(userID uint) (*UserInfoResult, error) {
	user, err := s.repo.User.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 查询用户角色
	roles, err := s.repo.Role.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	roleCodes := make([]string, 0, len(roles))
	roleIDs := make([]uint, 0, len(roles))
	for _, r := range roles {
		roleCodes = append(roleCodes, r.RoleCode)
		roleIDs = append(roleIDs, r.ID)
	}

	// 查询按钮权限（根据角色去重）
	buttons := make([]string, 0)
	if len(roleIDs) > 0 {
		perms, err := s.repo.Permission.FindByRoleIDs(roleIDs)
		if err != nil {
			return nil, err
		}
		seen := make(map[string]struct{})
		for _, p := range perms {
			if _, ok := seen[p.AuthMark]; !ok {
				seen[p.AuthMark] = struct{}{}
				buttons = append(buttons, p.AuthMark)
			}
		}
	}

	return &UserInfoResult{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Roles:    roleCodes,
		Buttons:  buttons,
	}, nil
}
