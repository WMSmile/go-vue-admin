package middleware

import (
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// whitelist paths that every authenticated user may reach.
var casbinWhitelist = map[string]bool{
	"/api/v1/menus/tree":  true,
	"/api/v1/auth/me":     true,
	"/api/v1/auth/logout": true,
}

// CasbinAuth enforces the request against role policies (dual backend check).
func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Request.URL.Path
		if casbinWhitelist[obj] {
			c.Next()
			return
		}

		rolesVal, ok := c.Get("roles")
		roleList, ok2 := rolesVal.([]string)
		if !ok || !ok2 || len(roleList) == 0 {
			utils.Fail(c, "forbidden")
			c.Abort()
			return
		}

		act := c.Request.Method
		allowed := false
		for _, role := range roleList {
			ok, err := global.Enforcer.Enforce(role, obj, act)
			if err == nil && ok {
				allowed = true
				break
			}
		}
		if !allowed {
			utils.Fail(c, "no permission")
			c.Abort()
			return
		}
		c.Next()
	}
}
