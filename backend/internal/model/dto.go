package model

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
