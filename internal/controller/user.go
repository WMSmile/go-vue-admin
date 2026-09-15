package controller

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserReq struct {
	ID       uint     `json:"id"`
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Status   int      `json:"status"`
	RoleIDs  []uint   `json:"roleIds"`
}

func ListUsers(c *gin.Context) {
	var users []models.User
	global.DB.Preload("Roles").Find(&users)
	utils.Ok(c, gin.H{"list": users, "total": len(users)})
}

func CreateUser(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	pwd := req.Password
	if pwd == "" {
		pwd = "123456"
	}
	hash, _ := utils.HashPassword(pwd)
	user := models.User{
		Username: req.Username,
		Password: hash,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}
	if err := global.DB.Create(&user).Error; err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		global.DB.Find(&roles, req.RoleIDs)
		global.DB.Model(&user).Association("Roles").Append(roles)
	}
	utils.Ok(c, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := global.DB.Preload("Roles").First(&user, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req UserReq
	c.ShouldBindJSON(&req)
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Status = req.Status
	if req.Password != "" {
		h, _ := utils.HashPassword(req.Password)
		user.Password = h
	}
	global.DB.Save(&user)
	if req.RoleIDs != nil {
		var roles []models.Role
		global.DB.Find(&roles, req.RoleIDs)
		global.DB.Model(&user).Association("Roles").Replace(roles)
	}
	utils.Ok(c, user)
}

func DeleteUser(c *gin.Context) {
	global.DB.Delete(&models.User{}, c.Param("id"))
	utils.Ok(c, nil)
}
