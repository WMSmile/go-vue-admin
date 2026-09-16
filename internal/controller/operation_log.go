package controller

import (
	"strconv"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// ListOperationLogs returns a paginated, filterable list of operation logs.
func ListOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	username := c.Query("username")
	method := c.Query("method")
	path := c.Query("path")

	db := global.DB.Model(&models.OperationLog{})
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if method != "" {
		db = db.Where("method = ?", method)
	}
	if path != "" {
		db = db.Where("path LIKE ?", "%"+path+"%")
	}

	var total int64
	db.Count(&total)

	var list []models.OperationLog
	db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	utils.Ok(c, gin.H{"list": list, "total": total})
}

// DeleteOperationLog removes a single log entry.
func DeleteOperationLog(c *gin.Context) {
	global.DB.Delete(&models.OperationLog{}, c.Param("id"))
	utils.Ok(c, nil)
}

// ClearOperationLogs removes all log entries.
func ClearOperationLogs(c *gin.Context) {
	global.DB.Where("1 = 1").Delete(&models.OperationLog{})
	utils.Ok(c, nil)
}
