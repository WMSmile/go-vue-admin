package controller

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// Dashboard returns simple overview statistics.
func Dashboard(c *gin.Context) {
	var users, roles, menus int64
	global.DB.Model(&models.User{}).Count(&users)
	global.DB.Model(&models.Role{}).Count(&roles)
	global.DB.Model(&models.Menu{}).Count(&menus)
	utils.Ok(c, gin.H{
		"users": users,
		"roles": roles,
		"menus": menus,
	})
}
