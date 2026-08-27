package model

import "time"

// Role 角色表
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoleCode    string    `gorm:"uniqueIndex;size:50;not null;comment:角色代码，如 R_SUPER" json:"roleCode"`
	RoleName    string    `gorm:"size:50;not null;comment:角色名称" json:"roleName"`
	Description string    `gorm:"size:255" json:"description"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`

	Menus       []*Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
