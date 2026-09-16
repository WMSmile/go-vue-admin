package initialize

import (
	"log"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"
)

// Seed creates the default admin/user, menu tree and casbin policies.
func Seed() {
	// roles
	var adminRole models.Role
	global.DB.Where(models.Role{Keyword: "admin"}).FirstOrCreate(&adminRole)
	adminRole.Name = "管理员"
	adminRole.Description = "超级管理员"
	adminRole.Status = 1
	global.DB.Save(&adminRole)

	var userRole models.Role
	global.DB.Where(models.Role{Keyword: "user"}).FirstOrCreate(&userRole)
	userRole.Name = "普通用户"
	userRole.Description = "普通用户"
	userRole.Status = 1
	global.DB.Save(&userRole)

	// admin user
	hash, _ := utils.HashPassword("admin123")
	var admin models.User
	global.DB.Where(models.User{Username: "admin"}).FirstOrCreate(&admin)
	admin.Password = hash
	admin.Nickname = "Admin"
	admin.Status = 1
	global.DB.Save(&admin)
	global.DB.Model(&admin).Association("Roles").Append(&adminRole)

	// menu tree
	createMenu := func(m models.Menu) models.Menu {
		var exist models.Menu
		if err := global.DB.Where(models.Menu{Name: m.Name}).First(&exist).Error; err != nil {
			global.DB.Create(&m)
			return m
		}
		// update in place so structural changes (e.g. ParentID/Path) apply on re-seed
		global.DB.Model(&exist).Updates(map[string]interface{}{
			"title": m.Title, "icon": m.Icon, "path": m.Path, "component": m.Component,
			"sort": m.Sort, "type": m.Type, "permission": m.Permission, "api": m.Api,
			"method": m.Method, "parent_id": m.ParentID, "status": m.Status,
		})
		m.ID = exist.ID
		return m
	}
	sys := createMenu(models.Menu{Name: "System", Title: "系统管理", Icon: "Setting", Path: "/system", Component: "Layout", Sort: 1, Type: models.MenuTypeCatalog})
	userMenu := createMenu(models.Menu{Name: "User", Title: "用户管理", Icon: "User", Path: "user", Component: "system/user/index", Sort: 1, Type: models.MenuTypeMenu, Api: "/api/v1/users", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "UserAdd", Title: "新增用户", Permission: "user:add", Api: "/api/v1/users", Method: "POST", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	createMenu(models.Menu{Name: "UserEdit", Title: "编辑用户", Permission: "user:edit", Api: "/api/v1/users/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	createMenu(models.Menu{Name: "UserDel", Title: "删除用户", Permission: "user:del", Api: "/api/v1/users/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	roleMenu := createMenu(models.Menu{Name: "Role", Title: "角色管理", Icon: "Avatar", Path: "role", Component: "system/role/index", Sort: 2, Type: models.MenuTypeMenu, Api: "/api/v1/roles", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "RoleAdd", Title: "新增角色", Permission: "role:add", Api: "/api/v1/roles", Method: "POST", Type: models.MenuTypeButton, ParentID: roleMenu.ID})
	createMenu(models.Menu{Name: "RoleEdit", Title: "编辑角色", Permission: "role:edit", Api: "/api/v1/roles/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: roleMenu.ID})
	createMenu(models.Menu{Name: "RoleDel", Title: "删除角色", Permission: "role:del", Api: "/api/v1/roles/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: roleMenu.ID})
	createMenu(models.Menu{Name: "RoleAssign", Title: "分配菜单", Permission: "role:assign", Api: "/api/v1/roles/menus", Method: "POST", Type: models.MenuTypeButton, ParentID: roleMenu.ID})
	menuMenu := createMenu(models.Menu{Name: "Menu", Title: "菜单管理", Icon: "Menu", Path: "menu", Component: "system/menu/index", Sort: 3, Type: models.MenuTypeMenu, Api: "/api/v1/menus", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "MenuAdd", Title: "新增菜单", Permission: "menu:add", Api: "/api/v1/menus", Method: "POST", Type: models.MenuTypeButton, ParentID: menuMenu.ID})
	createMenu(models.Menu{Name: "MenuEdit", Title: "编辑菜单", Permission: "menu:edit", Api: "/api/v1/menus/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: menuMenu.ID})
	createMenu(models.Menu{Name: "MenuDel", Title: "删除菜单", Permission: "menu:del", Api: "/api/v1/menus/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: menuMenu.ID})
	oplog := createMenu(models.Menu{Name: "OperationLog", Title: "操作日志", Icon: "Document", Path: "oplog", Component: "system/oplog/index", Sort: 4, Type: models.MenuTypeMenu, Api: "/api/v1/operation-logs", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "OpLogDel", Title: "删除日志", Permission: "oplog:del", Api: "/api/v1/operation-logs/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: oplog.ID})
	createMenu(models.Menu{Name: "OpLogClear", Title: "清空日志", Permission: "oplog:clear", Api: "/api/v1/operation-logs", Method: "DELETE", Type: models.MenuTypeButton, ParentID: oplog.ID})
	dict := createMenu(models.Menu{Name: "Dict", Title: "字典管理", Icon: "Collection", Path: "dict", Component: "system/dict/index", Sort: 5, Type: models.MenuTypeMenu, Api: "/api/v1/dict-types", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "DictAdd", Title: "新增字典", Permission: "dict:add", Api: "/api/v1/dict-types", Method: "POST", Type: models.MenuTypeButton, ParentID: dict.ID})
	createMenu(models.Menu{Name: "DictEdit", Title: "编辑字典", Permission: "dict:edit", Api: "/api/v1/dict-types/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: dict.ID})
	createMenu(models.Menu{Name: "DictDel", Title: "删除字典", Permission: "dict:del", Api: "/api/v1/dict-types/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: dict.ID})
	taskM := createMenu(models.Menu{Name: "Task", Title: "定时任务", Icon: "Timer", Path: "task", Component: "system/task/index", Sort: 6, Type: models.MenuTypeMenu, Api: "/api/v1/tasks", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "TaskAdd", Title: "新增任务", Permission: "task:add", Api: "/api/v1/tasks", Method: "POST", Type: models.MenuTypeButton, ParentID: taskM.ID})
	createMenu(models.Menu{Name: "TaskEdit", Title: "编辑任务", Permission: "task:edit", Api: "/api/v1/tasks/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: taskM.ID})
	createMenu(models.Menu{Name: "TaskDel", Title: "删除任务", Permission: "task:del", Api: "/api/v1/tasks/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: taskM.ID})
	createMenu(models.Menu{Name: "TaskRun", Title: "执行任务", Permission: "task:run", Api: "/api/v1/tasks/*/run", Method: "POST", Type: models.MenuTypeButton, ParentID: taskM.ID})
	configM := createMenu(models.Menu{Name: "Config", Title: "参数设置", Icon: "SetUp", Path: "config", Component: "system/config/index", Sort: 7, Type: models.MenuTypeMenu, Api: "/api/v1/configs", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "ConfigAdd", Title: "新增参数", Permission: "config:add", Api: "/api/v1/configs", Method: "POST", Type: models.MenuTypeButton, ParentID: configM.ID})
	createMenu(models.Menu{Name: "ConfigEdit", Title: "编辑参数", Permission: "config:edit", Api: "/api/v1/configs/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: configM.ID})
	createMenu(models.Menu{Name: "ConfigDel", Title: "删除参数", Permission: "config:del", Api: "/api/v1/configs/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: configM.ID})

	// seed default system parameters (branding etc.)
	seedConfig := func(key, value, name, group string, sort int) {
		var c models.SysConfig
		if err := global.DB.Where(models.SysConfig{Key: key}).First(&c).Error; err != nil {
			global.DB.Create(&models.SysConfig{Key: key, Value: value, Name: name, Group: group, Sort: sort, Status: 1})
		}
	}
	seedConfig("site.name", "Go Vue Admin", "系统名称", "基础设置", 1)
	seedConfig("site.loginTitle", "Go Vue Admin 管理系统", "登录页标题", "基础设置", 2)
	seedConfig("site.copyright", "© 2026 Go Vue Admin. All Rights Reserved.", "版权信息", "基础设置", 3)
	seedConfig("site.icp", "", "备案号", "基础设置", 4)

	// Dashboard: a top-level page shown directly after login
	dash := createMenu(models.Menu{Name: "Dashboard", Title: "仪表盘", Icon: "DataBoard", Path: "/dashboard", Component: "monitor/dashboard/index", Sort: 0, Type: models.MenuTypeMenu, Api: "/api/v1/dashboard", Method: "GET", ParentID: 0})
	// Monitor: a catalog container for other monitoring functions (logs, online users, etc.)
	mon := createMenu(models.Menu{Name: "Monitor", Title: "系统监控", Icon: "Monitor", Path: "/monitor", Component: "Layout", Sort: 3, Type: models.MenuTypeCatalog})

	// link menus to admin
	var allMenus []models.Menu
	global.DB.Find(&allMenus)
	global.DB.Model(&adminRole).Association("Menus").Append(allMenus)

	// link the Monitor catalog + Dashboard page to the normal user
	global.DB.Model(&userRole).Association("Menus").Append(&mon, &dash)

	// casbin policies
	addPolicy := func(sub, obj, act string) {
		if has, _ := global.Enforcer.HasPolicy(sub, obj, act); !has {
			global.Enforcer.AddPolicy(sub, obj, act)
		}
	}
	addPolicy("admin", "/*", "*")
	if dash.Api != "" && dash.Method != "" {
		addPolicy("user", dash.Api, dash.Method)
	}
	if err := global.Enforcer.SavePolicy(); err != nil {
		log.Println("save casbin policy:", err)
	}

	// demo API key for admin (machine/skill access). Generated once and reused.
	var demoKey models.ApiKey
	global.DB.Where(models.ApiKey{Name: "demo-skill"}).FirstOrCreate(&demoKey)
	if demoKey.Key == "" {
		if k, gerr := utils.GenerateAPIKey(); gerr == nil {
			demoKey.UserID = admin.ID
			demoKey.Key = k
			demoKey.Status = 1
			global.DB.Save(&demoKey)
		}
	}
	log.Printf("seed data ready | demo api key (admin): %s", demoKey.Key)
	log.Println("tip: call APIs with header  X-API-Key: <key>  (or Authorization: Bearer <token>)")
}
