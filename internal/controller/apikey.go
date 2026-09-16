package controller

import (
	"time"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateApiKeyReq struct {
	Name      string  `json:"name"`
	Scope     string  `json:"scope"`     // all | readonly
	ExpiresAt *string `json:"expiresAt"` // YYYY-MM-DD, optional (nil = never)
}

// apiKeyDTO avoids serializing the raw key and the owning user (with password hash).
type apiKeyDTO struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	MaskedKey  string     `json:"maskedKey"`
	Status     int        `json:"status"`
	Scope      string     `json:"scope"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	Expired    bool       `json:"expired"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func maskKey(k string) string {
	if len(k) <= 10 {
		return "****"
	}
	return k[:6] + "••••••••" + k[len(k)-4:]
}

// CreateApiKey issues a new API key for the current user. The raw key is returned
// only once (key/secret fields) and can be used by a skill/agent as the X-API-Key header.
func CreateApiKey(c *gin.Context) {
	uid, _ := c.Get("userId")

	var req CreateApiKeyReq
	c.ShouldBindJSON(&req)
	name := req.Name
	if name == "" {
		name = "skill-key"
	}

	scope := req.Scope
	if scope != "readonly" {
		scope = "all"
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err != nil {
			utils.Fail(c, "过期时间格式应为 YYYY-MM-DD")
			return
		}
		expiresAt = &t
	}

	key, err := utils.GenerateAPIKey()
	if err != nil {
		utils.Fail(c, "generate api key failed")
		return
	}

	ak := models.ApiKey{
		UserID:    uid.(uint),
		Name:      name,
		Key:       key,
		Status:    1,
		Scope:     scope,
		ExpiresAt: expiresAt,
	}
	if err := global.DB.Create(&ak).Error; err != nil {
		utils.Fail(c, "create failed")
		return
	}

	utils.Ok(c, gin.H{
		"id":        ak.ID,
		"name":      ak.Name,
		"key":       key,
		"secret":    key,
		"maskedKey": maskKey(key),
		"scope":     ak.Scope,
		"expiresAt": ak.ExpiresAt,
		"hint":      "请妥善保存，密钥仅展示一次",
		"createdAt": ak.CreatedAt,
	})
}

// ListApiKeys lists the current user's API keys (raw keys masked, user omitted).
func ListApiKeys(c *gin.Context) {
	uid, _ := c.Get("userId")
	var keys []models.ApiKey
	global.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&keys)

	now := time.Now()
	dtos := make([]apiKeyDTO, 0, len(keys))
	for _, k := range keys {
		expired := k.ExpiresAt != nil && k.ExpiresAt.Before(now)
		dtos = append(dtos, apiKeyDTO{
			ID:         k.ID,
			Name:       k.Name,
			MaskedKey:  maskKey(k.Key),
			Status:     k.Status,
			Scope:      k.Scope,
			ExpiresAt:  k.ExpiresAt,
			Expired:    expired,
			LastUsedAt: k.LastUsedAt,
			CreatedAt:  k.CreatedAt,
		})
	}
	utils.Ok(c, gin.H{"list": dtos})
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
