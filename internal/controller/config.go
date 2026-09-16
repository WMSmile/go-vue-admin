package controller

import (
	"strconv"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// ListConfigs returns a paginated, filterable list of system parameters.
func ListConfigs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	key := c.Query("key")
	name := c.Query("name")
	group := c.Query("group")

	db := global.DB.Model(&models.SysConfig{})
	if key != "" {
		db = db.Where("config_key LIKE ?", "%"+key+"%")
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if group != "" {
		db = db.Where("config_group = ?", group)
	}

	var total int64
	db.Count(&total)
	var list []models.SysConfig
	db.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	utils.Ok(c, gin.H{"list": list, "total": total})
}

// ConfigMap returns {key: value} for all enabled parameters. It is exposed
// publicly so the login page / footer can render branding without auth.
func ConfigMap(c *gin.Context) {
	var list []models.SysConfig
	global.DB.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&list)
	m := make(map[string]string, len(list))
	for _, cfg := range list {
		m[cfg.Key] = cfg.Value
	}
	utils.Ok(c, m)
}

type ConfigReq struct {
	Key    string `json:"key" binding:"required"`
	Value  string `json:"value"`
	Name   string `json:"name" binding:"required"`
	Group  string `json:"group"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// CreateConfig creates a new system parameter (key must be unique).
func CreateConfig(c *gin.Context) {
	var req ConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	var cnt int64
	global.DB.Model(&models.SysConfig{}).Where("config_key = ?", req.Key).Count(&cnt)
	if cnt > 0 {
		utils.Fail(c, "参数键已存在")
		return
	}
	cfg := models.SysConfig{
		Key: req.Key, Value: req.Value, Name: req.Name,
		Group: req.Group, Sort: req.Sort, Status: req.Status, Remark: req.Remark,
	}
	global.DB.Create(&cfg)
	utils.Ok(c, cfg)
}

// UpdateConfig updates an existing system parameter. The key is immutable.
func UpdateConfig(c *gin.Context) {
	var cfg models.SysConfig
	if err := global.DB.First(&cfg, c.Param("id")).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req ConfigReq
	c.ShouldBindJSON(&req)
	cfg.Value = req.Value
	cfg.Name = req.Name
	cfg.Group = req.Group
	cfg.Sort = req.Sort
	cfg.Status = req.Status
	cfg.Remark = req.Remark
	global.DB.Save(&cfg)
	utils.Ok(c, cfg)
}

// DeleteConfig removes a single system parameter.
func DeleteConfig(c *gin.Context) {
	global.DB.Delete(&models.SysConfig{}, c.Param("id"))
	utils.Ok(c, nil)
}
