package main

import "golang-admin-best/internal/model"

// seedMenu 菜单种子结构（镜像前端 src/router/modules 的 AppRouteRecord）
type seedMenu struct {
	Path          string
	Name          string
	Component     string
	Redirect      string
	Title         string
	Icon          string
	KeepAlive     bool
	IsHide        bool
	IsHideTab     bool
	IsIframe      bool
	IsFullPage    bool
	FixedTab      bool
	Link          string
	ActivePath    string
	ShowTextBadge string
	Roles         []string         // 可见角色；为空表示继承父级（默认全部）
	Auth          []model.AuthItem // 按钮级权限（authList）
	Children      []seedMenu
}

// menuSeeds 返回完整菜单树，顺序与前端 router/modules/index.ts 一致
func menuSeeds() []seedMenu {
	return []seedMenu{
		dashboardSeed(),
		templateSeed(),
		widgetsSeed(),
		examplesSeed(),
		systemSeed(),
		articleSeed(),
		resultSeed(),
		exceptionSeed(),
		safeguardSeed(),
		// help 是并列的多个顶级项
		helpDocumentSeed(),
		helpLiteVersionSeed(),
		helpOldVersionSeed(),
		helpChangeLogSeed(),
	}
}

func dashboardSeed() seedMenu {
	return seedMenu{
		Path: "/dashboard", Name: "Dashboard", Component: "/index/index",
		Title: "menus.dashboard.title", Icon: "ri:pie-chart-line",
		Roles: []string{"R_SUPER", "R_ADMIN"},
		Children: []seedMenu{
			{Path: "console", Name: "Console", Component: "/dashboard/console", Title: "menus.dashboard.console", Icon: "ri:home-smile-2-line", FixedTab: true},
			{Path: "analysis", Name: "Analysis", Component: "/dashboard/analysis", Title: "menus.dashboard.analysis", Icon: "ri:align-item-bottom-line"},
			{Path: "ecommerce", Name: "Ecommerce", Component: "/dashboard/ecommerce", Title: "menus.dashboard.ecommerce", Icon: "ri:bar-chart-box-line"},
		},
	}
}

func templateSeed() seedMenu {
	return seedMenu{
		Path: "/template", Name: "Template", Component: "/index/index",
		Title: "menus.template.title", Icon: "ri:apps-2-line",
		Children: []seedMenu{
			{Path: "cards", Name: "Cards", Component: "/template/cards", Title: "menus.template.cards", Icon: "ri:wallet-line"},
			{Path: "banners", Name: "Banners", Component: "/template/banners", Title: "menus.template.banners", Icon: "ri:rectangle-line"},
			{Path: "charts", Name: "Charts", Component: "/template/charts", Title: "menus.template.charts", Icon: "ri:bar-chart-box-line"},
			{Path: "map", Name: "Map", Component: "/template/map", Title: "menus.template.map", Icon: "ri:map-pin-line", KeepAlive: true},
			{Path: "chat", Name: "Chat", Component: "/template/chat", Title: "menus.template.chat", Icon: "ri:message-3-line", KeepAlive: true},
			{Path: "calendar", Name: "Calendar", Component: "/template/calendar", Title: "menus.template.calendar", Icon: "ri:calendar-2-line", KeepAlive: true},
			{Path: "pricing", Name: "Pricing", Component: "/template/pricing", Title: "menus.template.pricing", Icon: "ri:money-cny-box-line", KeepAlive: true, IsFullPage: true},
		},
	}
}

func widgetsSeed() seedMenu {
	return seedMenu{
		Path: "/widgets", Name: "Widgets", Component: "/index/index",
		Title: "menus.widgets.title", Icon: "ri:apps-2-add-line",
		Children: []seedMenu{
			{Path: "icon", Name: "Icon", Component: "/widgets/icon", Title: "menus.widgets.icon", Icon: "ri:palette-line", KeepAlive: true},
			{Path: "image-crop", Name: "ImageCrop", Component: "/widgets/image-crop", Title: "menus.widgets.imageCrop", Icon: "ri:screenshot-line", KeepAlive: true},
			{Path: "excel", Name: "Excel", Component: "/widgets/excel", Title: "menus.widgets.excel", Icon: "ri:download-2-line", KeepAlive: true},
			{Path: "video", Name: "Video", Component: "/widgets/video", Title: "menus.widgets.video", Icon: "ri:vidicon-line", KeepAlive: true},
			{Path: "count-to", Name: "CountTo", Component: "/widgets/count-to", Title: "menus.widgets.countTo", Icon: "ri:anthropic-line"},
			{Path: "wang-editor", Name: "WangEditor", Component: "/widgets/wang-editor", Title: "menus.widgets.wangEditor", Icon: "ri:t-box-line", KeepAlive: true},
			{Path: "watermark", Name: "Watermark", Component: "/widgets/watermark", Title: "menus.widgets.watermark", Icon: "ri:water-flash-line", KeepAlive: true},
			{Path: "context-menu", Name: "ContextMenu", Component: "/widgets/context-menu", Title: "menus.widgets.contextMenu", Icon: "ri:menu-2-line", KeepAlive: true},
			{Path: "qrcode", Name: "Qrcode", Component: "/widgets/qrcode", Title: "menus.widgets.qrcode", Icon: "ri:qr-code-line", KeepAlive: true},
			{Path: "drag", Name: "Drag", Component: "/widgets/drag", Title: "menus.widgets.drag", Icon: "ri:drag-move-fill", KeepAlive: true},
			{Path: "text-scroll", Name: "TextScroll", Component: "/widgets/text-scroll", Title: "menus.widgets.textScroll", Icon: "ri:input-method-line", KeepAlive: true},
			{Path: "fireworks", Name: "Fireworks", Component: "/widgets/fireworks", Title: "menus.widgets.fireworks", Icon: "ri:magic-line", KeepAlive: true, ShowTextBadge: "Hot"},
			{Path: "/outside/iframe/elementui", Name: "ElementUI", Title: "menus.widgets.elementUI", Icon: "ri:apps-2-line", IsIframe: true, Link: "https://element-plus.org/zh-CN/component/overview.html"},
		},
	}
}

func examplesSeed() seedMenu {
	return seedMenu{
		Path: "/examples", Name: "Examples", Component: "/index/index",
		Title: "menus.examples.title", Icon: "ri:sparkling-line",
		Children: []seedMenu{
			{
				Path: "permission", Name: "Permission", Title: "menus.examples.permission.title", Icon: "ri:fingerprint-line",
				Children: []seedMenu{
					{Path: "switch-role", Name: "PermissionSwitchRole", Component: "/examples/permission/switch-role", Title: "menus.examples.permission.switchRole", Icon: "ri:contacts-line", KeepAlive: true},
					{Path: "button-auth", Name: "PermissionButtonAuth", Component: "/examples/permission/button-auth", Title: "menus.examples.permission.buttonAuth", Icon: "ri:mouse-line", KeepAlive: true,
						Auth: []model.AuthItem{
							{Title: "新增", AuthMark: "add"}, {Title: "编辑", AuthMark: "edit"}, {Title: "删除", AuthMark: "delete"}, {Title: "导出", AuthMark: "export"},
							{Title: "查看", AuthMark: "view"}, {Title: "发布", AuthMark: "publish"}, {Title: "配置", AuthMark: "config"}, {Title: "管理", AuthMark: "manage"},
						}},
					{Path: "page-visibility", Name: "PermissionPageVisibility", Component: "/examples/permission/page-visibility", Title: "menus.examples.permission.pageVisibility", Icon: "ri:user-3-line", KeepAlive: true, Roles: []string{"R_SUPER"}},
				},
			},
			{Path: "tabs", Name: "Tabs", Component: "/examples/tabs", Title: "menus.examples.tabs", Icon: "ri:price-tag-line"},
			{Path: "tables/basic", Name: "TablesBasic", Component: "/examples/tables/basic", Title: "menus.examples.tablesBasic", Icon: "ri:layout-grid-line", KeepAlive: true},
			{Path: "tables", Name: "Tables", Component: "/examples/tables", Title: "menus.examples.tables", Icon: "ri:table-3", KeepAlive: true},
			{Path: "forms", Name: "Forms", Component: "/examples/forms", Title: "menus.examples.forms", Icon: "ri:table-view", KeepAlive: true},
			{Path: "form/search-bar", Name: "SearchBar", Component: "/examples/forms/search-bar", Title: "menus.examples.searchBar", Icon: "ri:table-line", KeepAlive: true},
			{Path: "tables/tree", Name: "TablesTree", Component: "/examples/tables/tree", Title: "menus.examples.tablesTree", Icon: "ri:layout-2-line", KeepAlive: true},
			{Path: "socket-chat", Name: "SocketChat", Component: "/examples/socket-chat", Title: "menus.examples.socketChat", Icon: "ri:shake-hands-line", KeepAlive: true},
		},
	}
}

func systemSeed() seedMenu {
	return seedMenu{
		Path: "/system", Name: "System", Component: "/index/index",
		Title: "menus.system.title", Icon: "ri:user-3-line",
		Roles: []string{"R_SUPER", "R_ADMIN"},
		Children: []seedMenu{
			{Path: "user", Name: "User", Component: "/system/user", Title: "menus.system.user", Icon: "ri:user-line", KeepAlive: true, Roles: []string{"R_SUPER", "R_ADMIN"}},
			{Path: "role", Name: "Role", Component: "/system/role", Title: "menus.system.role", Icon: "ri:user-settings-line", KeepAlive: true, Roles: []string{"R_SUPER"}},
			{Path: "user-center", Name: "UserCenter", Component: "/system/user-center", Title: "menus.system.userCenter", Icon: "ri:user-line", IsHide: true, KeepAlive: true, IsHideTab: true},
			{Path: "menu", Name: "Menus", Component: "/system/menu", Title: "menus.system.menu", Icon: "ri:menu-line", KeepAlive: true, Roles: []string{"R_SUPER"},
				Auth: []model.AuthItem{{Title: "新增", AuthMark: "add"}, {Title: "编辑", AuthMark: "edit"}, {Title: "删除", AuthMark: "delete"}}},
			{
				Path: "nested", Name: "Nested", Title: "menus.system.nested", Icon: "ri:menu-unfold-3-line", KeepAlive: true,
				Children: []seedMenu{
					{Path: "menu1", Name: "NestedMenu1", Component: "/system/nested/menu1", Title: "menus.system.menu1", Icon: "ri:align-justify", KeepAlive: true},
					{
						Path: "menu2", Name: "NestedMenu2", Title: "menus.system.menu2", Icon: "ri:align-justify", KeepAlive: true,
						Children: []seedMenu{
							{Path: "menu2-1", Name: "NestedMenu2-1", Component: "/system/nested/menu2", Title: "menus.system.menu21", Icon: "ri:align-justify", KeepAlive: true},
						},
					},
					{
						Path: "menu3", Name: "NestedMenu3", Title: "menus.system.menu3", Icon: "ri:align-justify", KeepAlive: true,
						Children: []seedMenu{
							{Path: "menu3-1", Name: "NestedMenu3-1", Component: "/system/nested/menu3", Title: "menus.system.menu31", KeepAlive: true},
							{
								Path: "menu3-2", Name: "NestedMenu3-2", Title: "menus.system.menu32", KeepAlive: true,
								Children: []seedMenu{
									{Path: "menu3-2-1", Name: "NestedMenu3-2-1", Component: "/system/nested/menu3/menu3-2", Title: "menus.system.menu321", KeepAlive: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func articleSeed() seedMenu {
	return seedMenu{
		Path: "/article", Name: "Article", Component: "/index/index",
		Title: "menus.article.title", Icon: "ri:book-2-line",
		Roles: []string{"R_SUPER", "R_ADMIN"},
		Children: []seedMenu{
			{Path: "article-list", Name: "ArticleList", Component: "/article/list", Title: "menus.article.articleList", Icon: "ri:article-line", KeepAlive: true,
				Auth: []model.AuthItem{{Title: "新增", AuthMark: "add"}, {Title: "编辑", AuthMark: "edit"}}},
			{Path: "detail/:id", Name: "ArticleDetail", Component: "/article/detail", Title: "menus.article.articleDetail", IsHide: true, KeepAlive: true, ActivePath: "/article/article-list"},
			{Path: "comment", Name: "ArticleComment", Component: "/article/comment", Title: "menus.article.comment", Icon: "ri:mail-line", KeepAlive: true},
			{Path: "publish", Name: "ArticlePublish", Component: "/article/publish", Title: "menus.article.articlePublish", Icon: "ri:telegram-2-line", KeepAlive: true,
				Auth: []model.AuthItem{{Title: "发布", AuthMark: "add"}}},
		},
	}
}

func resultSeed() seedMenu {
	return seedMenu{
		Path: "/result", Name: "Result", Component: "/index/index",
		Title: "menus.result.title", Icon: "ri:checkbox-circle-line",
		Children: []seedMenu{
			{Path: "success", Name: "ResultSuccess", Component: "/result/success", Title: "menus.result.success", Icon: "ri:checkbox-circle-line", KeepAlive: true},
			{Path: "fail", Name: "ResultFail", Component: "/result/fail", Title: "menus.result.fail", Icon: "ri:close-circle-line", KeepAlive: true},
		},
	}
}

func exceptionSeed() seedMenu {
	return seedMenu{
		Path: "/exception", Name: "Exception", Component: "/index/index",
		Title: "menus.exception.title", Icon: "ri:error-warning-line",
		Children: []seedMenu{
			{Path: "403", Name: "Exception403", Component: "/exception/403", Title: "menus.exception.forbidden", KeepAlive: true, IsHideTab: true, IsFullPage: true},
			{Path: "404", Name: "Exception404", Component: "/exception/404", Title: "menus.exception.notFound", KeepAlive: true, IsHideTab: true, IsFullPage: true},
			{Path: "500", Name: "Exception500", Component: "/exception/500", Title: "menus.exception.serverError", KeepAlive: true, IsHideTab: true, IsFullPage: true},
		},
	}
}

func safeguardSeed() seedMenu {
	return seedMenu{
		Path: "/safeguard", Name: "Safeguard", Component: "/index/index",
		Title: "menus.safeguard.title", Icon: "ri:shield-check-line",
		Children: []seedMenu{
			{Path: "server", Name: "SafeguardServer", Component: "/safeguard/server", Title: "menus.safeguard.server", Icon: "ri:hard-drive-3-line", KeepAlive: true},
		},
	}
}

// help 模块在前端是并列的 4 个顶级项（非父子结构）

func helpDocumentSeed() seedMenu {
	return seedMenu{Name: "Document", Title: "menus.help.document", Icon: "ri:bill-line", Link: "https://www.artd.pro/docs/zh/"}
}

func helpLiteVersionSeed() seedMenu {
	return seedMenu{Name: "LiteVersion", Title: "menus.help.liteVersion", Icon: "ri:bus-2-line", Link: "https://www.artd.pro/docs/zh/guide/lite-version.html"}
}

func helpOldVersionSeed() seedMenu {
	return seedMenu{Name: "OldVersion", Title: "menus.help.oldVersion", Icon: "ri:subway-line", Link: "https://www.artd.pro/v2/"}
}

func helpChangeLogSeed() seedMenu {
	return seedMenu{Path: "/change/log", Name: "ChangeLog", Component: "/change/log", Title: "menus.plan.log", Icon: "ri:gamepad-line", ShowTextBadge: "v3.0.2"}
}
