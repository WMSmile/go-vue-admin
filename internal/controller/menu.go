package controller

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type MenuReq struct {
	ID         uint   `json:"id"`
	ParentID   uint   `json:"parentId"`
	Name       string `json:"name"`
	Title      string `json:"title" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Sort       int    `json:"sort"`
	Type       int    `json:"type"`
	Permission string `json:"permission"`
	Api        string `json:"api"`
	Method     string `json:"method"`
	Status     int    `json:"status"`
	Hidden     int    `json:"hidden"`
}

func ListMenus(c *gin.Context) {
	var menus []models.Menu
	global.DB.Order("sort asc").Find(&menus)
	utils.Ok(c, menus)
}

func CreateMenu(c *gin.Context) {
	var req MenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	m := models.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Sort:       req.Sort,
		Type:       req.Type,
		Permission: req.Permission,
		Api:        req.Api,
		Method:     req.Method,
		Status:     1,
		Hidden:     req.Hidden,
	}
	global.DB.Create(&m)
	utils.Ok(c, m)
}

func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var m models.Menu
	if err := global.DB.First(&m, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req MenuReq
	c.ShouldBindJSON(&req)
	m.ParentID = req.ParentID
	m.Name = req.Name
	m.Title = req.Title
	m.Icon = req.Icon
	m.Path = req.Path
	m.Component = req.Component
	m.Sort = req.Sort
	m.Type = req.Type
	m.Permission = req.Permission
	m.Api = req.Api
	m.Method = req.Method
	m.Status = req.Status
	m.Hidden = req.Hidden
	global.DB.Save(&m)
	utils.Ok(c, m)
}

func DeleteMenu(c *gin.Context) {
	global.DB.Delete(&models.Menu{}, c.Param("id"))
	utils.Ok(c, nil)
}

// MenuTree returns the current user's menu tree (dynamic sidebar source).
func MenuTree(c *gin.Context) {
	rolesVal, _ := c.Get("roles")
	keywords, _ := rolesVal.([]string)
	tree := toTree(menusForRoles(keywords), 0)
	utils.Ok(c, tree)
}
