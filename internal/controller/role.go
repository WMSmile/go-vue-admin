package controller

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type RoleReq struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Keyword     string `json:"keyword" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type AssignMenusReq struct {
	RoleID  uint   `json:"roleId" binding:"required"`
	MenuIDs []uint `json:"menuIds"`
}

func ListRoles(c *gin.Context) {
	var roles []models.Role
	global.DB.Preload("Menus").Find(&roles)
	utils.Ok(c, gin.H{"list": roles, "total": len(roles)})
}

func CreateRole(c *gin.Context) {
	var req RoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	role := models.Role{Name: req.Name, Keyword: req.Keyword, Description: req.Description, Status: 1}
	global.DB.Create(&role)
	utils.Ok(c, role)
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := global.DB.First(&role, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req RoleReq
	c.ShouldBindJSON(&req)
	role.Name = req.Name
	role.Keyword = req.Keyword
	role.Description = req.Description
	role.Status = req.Status
	global.DB.Save(&role)
	utils.Ok(c, role)
}

func DeleteRole(c *gin.Context) {
	global.DB.Delete(&models.Role{}, c.Param("id"))
	utils.Ok(c, nil)
}

// AssignMenus links menus to a role and syncs casbin api policies (non-admin).
func AssignMenus(c *gin.Context) {
	var req AssignMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	var role models.Role
	if err := global.DB.First(&role, req.RoleID).Error; err != nil {
		utils.Fail(c, "role not found")
		return
	}
	var menus []models.Menu
	if len(req.MenuIDs) > 0 {
		global.DB.Find(&menus, req.MenuIDs)
	}
	global.DB.Model(&role).Association("Menus").Replace(menus)

	if role.Keyword != "admin" {
		global.Enforcer.RemoveFilteredPolicy(0, role.Keyword)
		for _, m := range menus {
			if m.Api != "" && m.Method != "" {
				global.Enforcer.AddPolicy(role.Keyword, m.Api, m.Method)
			}
		}
		global.Enforcer.SavePolicy()
	}
	utils.Ok(c, nil)
}
