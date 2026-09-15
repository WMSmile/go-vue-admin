package controller

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateApiKeyReq struct {
	Name string `json:"name"`
}

// CreateApiKey issues a new API key for the current user. The raw key is returned
// only once (secret field) and can be used by a skill/agent as the X-API-Key header.
func CreateApiKey(c *gin.Context) {
	uid, _ := c.Get("userId")

	var req CreateApiKeyReq
	c.ShouldBindJSON(&req)
	name := req.Name
	if name == "" {
		name = "skill-key"
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		utils.Fail(c, "generate api key failed")
		return
	}

	ak := models.ApiKey{UserID: uid.(uint), Name: name, Key: key, Status: 1}
	if err := global.DB.Create(&ak).Error; err != nil {
		utils.Fail(c, "create failed")
		return
	}

	utils.Ok(c, gin.H{
		"id":        ak.ID,
		"name":      ak.Name,
		"key":       key,
		"secret":    key,
		"hint":      "请妥善保存，密钥仅展示一次",
		"createdAt": ak.CreatedAt,
	})
}

// ListApiKeys lists the current user's API keys (raw keys omitted).
func ListApiKeys(c *gin.Context) {
	uid, _ := c.Get("userId")
	var keys []models.ApiKey
	global.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&keys)
	utils.Ok(c, keys)
}

// DeleteApiKey revokes a single API key owned by the current user.
func DeleteApiKey(c *gin.Context) {
	uid, _ := c.Get("userId")
	res := global.DB.Where("id = ? AND user_id = ?", c.Param("id"), uid).Delete(&models.ApiKey{})
	if res.RowsAffected == 0 {
		utils.Fail(c, "not found")
		return
	}
	utils.Ok(c, nil)
}
