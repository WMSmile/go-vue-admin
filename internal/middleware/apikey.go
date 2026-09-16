package middleware

import (
	"net/http"
	"time"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth authenticates a request via the "X-API-Key" header, and falls back
// to the JWT bearer token when no key is present. On success it loads the owning
// user's roles into the context so the subsequent Casbin middleware can enforce
// per-role permissions (query/create/update/delete) for machine callers (skills).
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			JWTAuth()(c)
			return
		}

		var ak models.ApiKey
		if err := global.DB.Preload("User.Roles").
			Where("`key` = ? AND status = ?", key, 1).First(&ak).Error; err != nil || ak.User.ID == 0 {
			utils.Unauthorized(c, "invalid api key")
			c.Abort()
			return
		}

		// expiry check
		if ak.ExpiresAt != nil && ak.ExpiresAt.Before(time.Now()) {
			utils.Forbidden(c, "api key expired")
			c.Abort()
			return
		}

		// scope check: readonly keys may only perform safe (read) methods
		if ak.Scope == "readonly" {
			switch c.Request.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
			default:
				utils.Forbidden(c, "该密钥仅允许只读操作 (GET)")
				c.Abort()
				return
			}
		}

		roles := make([]string, 0, len(ak.User.Roles))
		for _, r := range ak.User.Roles {
			roles = append(roles, r.Keyword)
		}

		now := time.Now()
		global.DB.Model(&ak).Update("last_used_at", &now)

		c.Set("userId", ak.User.ID)
		c.Set("username", ak.User.Username)
		c.Set("roles", roles)
		c.Set("apiKeyId", ak.ID)
		c.Next()
	}
}
