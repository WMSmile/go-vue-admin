package middleware

import (
	"strings"

	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth validates the bearer token and stores identity in context.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.Unauthorized(c, "missing token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
