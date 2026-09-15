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
	sys := createMenu(models.Menu{Name: "System", Title: "系统管理", Icon: "setting", Path: "/system", Component: "Layout", Sort: 1, Type: models.MenuTypeCatalog})
	userMenu := createMenu(models.Menu{Name: "User", Title: "用户管理", Icon: "user", Path: "user", Component: "system/user/index", Sort: 1, Type: models.MenuTypeMenu, Api: "/api/v1/users", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "UserAdd", Title: "新增用户", Permission: "user:add", Api: "/api/v1/users", Method: "POST", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	createMenu(models.Menu{Name: "UserEdit", Title: "编辑用户", Permission: "user:edit", Api: "/api/v1/users/*", Method: "PUT", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	createMenu(models.Menu{Name: "UserDel", Title: "删除用户", Permission: "user:del", Api: "/api/v1/users/*", Method: "DELETE", Type: models.MenuTypeButton, ParentID: userMenu.ID})
	createMenu(models.Menu{Name: "Role", Title: "角色管理", Icon: "role", Path: "role", Component: "system/role/index", Sort: 2, Type: models.MenuTypeMenu, Api: "/api/v1/roles", Method: "GET", ParentID: sys.ID})
	createMenu(models.Menu{Name: "Menu", Title: "菜单管理", Icon: "menu", Path: "menu", Component: "system/menu/index", Sort: 3, Type: models.MenuTypeMenu, Api: "/api/v1/menus", Method: "GET", ParentID: sys.ID})

	// Dashboard: a top-level page shown directly after login
	dash := createMenu(models.Menu{Name: "Dashboard", Title: "仪表盘", Icon: "dashboard", Path: "/dashboard", Component: "monitor/dashboard/index", Sort: 0, Type: models.MenuTypeMenu, Api: "/api/v1/dashboard", Method: "GET", ParentID: 0})
	// Monitor: a catalog container for other monitoring functions (logs, online users, etc.)
	mon := createMenu(models.Menu{Name: "Monitor", Title: "系统监控", Icon: "monitor", Path: "/monitor", Component: "Layout", Sort: 3, Type: models.MenuTypeCatalog})

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
