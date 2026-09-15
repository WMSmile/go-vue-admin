package controller

import (
	"sort"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}

	var user models.User
	if err := global.DB.Preload("Roles.Menus").Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, "user not found")
		return
	}
	if user.Status != 1 {
		utils.Fail(c, "user disabled")
		return
	}
	if !utils.CheckPassword(user.Password, req.Password) {
		utils.Fail(c, "wrong password")
		return
	}

	roleKeywords := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roleKeywords = append(roleKeywords, r.Keyword)
	}

	token, err := utils.GenerateToken(user.ID, user.Username, roleKeywords)
	if err != nil {
		utils.Fail(c, "token generation failed")
		return
	}

	utils.Ok(c, gin.H{
		"token":       token,
		"user":        user,
		"menus":       buildMenuTree(user.Roles),
		"permissions": buildPermissions(user.Roles),
	})
}

func Me(c *gin.Context) {
	uid, _ := c.Get("userId")
	var user models.User
	if err := global.DB.Preload("Roles").First(&user, uid).Error; err != nil {
		utils.Fail(c, "user not found")
		return
	}
	utils.Ok(c, user)
}

func Logout(c *gin.Context) {
	utils.Ok(c, nil)
}

// ----- menu tree / permission helpers -----

func buildMenuTree(roles []models.Role) []models.Menu {
	all := collectRoleMenus(roles, models.MenuTypeButton)
	return toTree(all, 0)
}

func buildPermissions(roles []models.Role) []string {
	perms := make([]string, 0)
	seen := map[string]bool{}
	for _, role := range roles {
		for _, m := range role.Menus {
			if m.Type == models.MenuTypeButton && m.Permission != "" && !seen[m.Permission] {
				seen[m.Permission] = true
				perms = append(perms, m.Permission)
			}
		}
	}
	return perms
}

func menusForRoles(keywords []string) []models.Menu {
	var roles []models.Role
	global.DB.Preload("Menus").Where("keyword IN ?", keywords).Find(&roles)
	return collectRoleMenus(roles, models.MenuTypeButton)
}

func collectRoleMenus(roles []models.Role, excludeType int) []models.Menu {
	all := make([]models.Menu, 0)
	seen := map[uint]bool{}
	for _, role := range roles {
		for _, m := range role.Menus {
			if m.Type == excludeType || seen[m.ID] {
				continue
			}
			seen[m.ID] = true
			all = append(all, m)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Sort < all[j].Sort })
	return all
}

func toTree(menus []models.Menu, parentID uint) []models.Menu {
	var res []models.Menu
	for _, m := range menus {
		if m.ParentID == parentID {
			m.Children = toTree(menus, m.ID)
			res = append(res, m)
		}
	}
	return res
}
