package model

import "time"

// Menu 菜单表
type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0;index;comment:父菜单ID，0表示一级菜单" json:"parentId"`
	Path      string    `gorm:"size:255;not null;comment:路由路径" json:"path"`
	Name      string    `gorm:"size:100;not null;comment:路由名称" json:"name"`
	Component string    `gorm:"size:255;comment:组件路径" json:"component"`
	Redirect  string    `gorm:"size:255;comment:重定向路径" json:"redirect"`
	Title     string    `gorm:"size:100;not null;comment:菜单标题" json:"title"`
	Icon      string    `gorm:"size:100;comment:图标" json:"icon"`
	KeepAlive int8      `gorm:"default:0;comment:是否缓存" json:"keepAlive"`
	IsHide    int8      `gorm:"default:0;comment:是否隐藏" json:"isHide"`
	IsIframe  int8      `gorm:"default:0;comment:是否iframe" json:"isIframe"`
	Link      string    `gorm:"size:255;comment:外部链接" json:"link"`
	Sort      int       `gorm:"default:0;comment:排序" json:"sort"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"createdAt"`

	Permissions []*Permission `gorm:"foreignKey:MenuID" json:"permissions,omitempty"`
}

func (Menu) TableName() string {
	return "menus"
}
