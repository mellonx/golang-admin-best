package model

// Permission 权限表（按钮级权限）
type Permission struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	MenuID   uint   `gorm:"index;not null;comment:所属菜单" json:"menuId"`
	Title    string `gorm:"size:50;not null;comment:权限名称" json:"title"`
	AuthMark string `gorm:"size:50;not null;comment:权限标识，如 add, edit, delete" json:"authMark"`

	Roles []*Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}
