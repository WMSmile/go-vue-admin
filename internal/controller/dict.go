package controller

import (
	"strconv"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// ---------- DictType ----------

// ListDictTypes returns a paginated, filterable list of dictionary types.
func ListDictTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	code := c.Query("code")
	name := c.Query("name")

	db := global.DB.Model(&models.DictType{})
	if code != "" {
		db = db.Where("code LIKE ?", "%"+code+"%")
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	db.Count(&total)
	var list []models.DictType
	db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	utils.Ok(c, gin.H{"list": list, "total": total})
}

type DictTypeReq struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// CreateDictType creates a new dictionary type (code must be unique).
func CreateDictType(c *gin.Context) {
	var req DictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	var cnt int64
	global.DB.Model(&models.DictType{}).Where("code = ?", req.Code).Count(&cnt)
	if cnt > 0 {
		utils.Fail(c, "字典编码已存在")
		return
	}
	dt := models.DictType{Code: req.Code, Name: req.Name, Status: req.Status, Remark: req.Remark}
	global.DB.Create(&dt)
	utils.Ok(c, dt)
}

// UpdateDictType updates an existing dictionary type.
func UpdateDictType(c *gin.Context) {
	var dt models.DictType
	if err := global.DB.First(&dt, c.Param("id")).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req DictTypeReq
	c.ShouldBindJSON(&req)
	dt.Code = req.Code
	dt.Name = req.Name
	dt.Status = req.Status
	dt.Remark = req.Remark
	global.DB.Save(&dt)
	utils.Ok(c, dt)
}

// DeleteDictType removes a type and all its dictionary entries (cascade).
func DeleteDictType(c *gin.Context) {
	global.DB.Where("type_id = ?", c.Param("id")).Delete(&models.DictData{})
	global.DB.Delete(&models.DictType{}, c.Param("id"))
	utils.Ok(c, nil)
}

// ---------- DictData ----------

// ListDictData returns dictionary entries, optionally filtered by type.
func ListDictData(c *gin.Context) {
	typeID := c.Query("typeId")
	typeCode := c.Query("typeCode")
	label := c.Query("label")

	db := global.DB.Model(&models.DictData{})
	if typeID != "" {
		db = db.Where("type_id = ?", typeID)
	}
	if typeCode != "" {
		db = db.Where("type_code = ?", typeCode)
	}
	if label != "" {
		db = db.Where("label LIKE ?", "%"+label+"%")
	}

	var list []models.DictData
	db.Order("sort ASC, id ASC").Find(&list)
	utils.Ok(c, gin.H{"list": list, "total": len(list)})
}

type DictDataReq struct {
	TypeID uint   `json:"typeId" binding:"required"`
	Label  string `json:"label" binding:"required"`
	Value  string `json:"value" binding:"required"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// CreateDictData creates a dictionary entry under a type.
func CreateDictData(c *gin.Context) {
	var req DictDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	var dt models.DictType
	if err := global.DB.First(&dt, req.TypeID).Error; err != nil {
		utils.Fail(c, "字典类型不存在")
		return
	}
	dd := models.DictData{
		TypeID:   req.TypeID,
		TypeCode: dt.Code,
		Label:    req.Label,
		Value:    req.Value,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	global.DB.Create(&dd)
	utils.Ok(c, dd)
}

// UpdateDictData updates a dictionary entry.
func UpdateDictData(c *gin.Context) {
	var dd models.DictData
	if err := global.DB.First(&dd, c.Param("id")).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req DictDataReq
	c.ShouldBindJSON(&req)
	dd.Label = req.Label
	dd.Value = req.Value
	dd.Sort = req.Sort
	dd.Status = req.Status
	dd.Remark = req.Remark
	global.DB.Save(&dd)
	utils.Ok(c, dd)
}

// DeleteDictData removes a single dictionary entry.
func DeleteDictData(c *gin.Context) {
	global.DB.Delete(&models.DictData{}, c.Param("id"))
	utils.Ok(c, nil)
}
