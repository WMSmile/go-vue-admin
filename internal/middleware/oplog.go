package middleware

import (
	"strings"
	"time"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"

	"github.com/gin-gonic/gin"
)

// skipPaths are framework/high-frequency endpoints we do not want to flood the
// operation log with (menus, current user, dashboard, and the log listing itself).
var skipLogPaths = map[string]bool{
	"/api/v1/menus/tree": true,
	"/api/v1/auth/me":    true,
	"/api/v1/dashboard":  true,
}

// OperationLog records every authenticated API request for audit purposes.
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if c.Request.Method == "OPTIONS" {
			return
		}
		path := c.Request.URL.Path
		if skipLogPaths[path] || strings.HasPrefix(path, "/api/v1/operation-logs") {
			return
		}

		var userID uint
		if v, ok := c.Get("userId"); ok {
			if u, ok := v.(uint); ok {
				userID = u
			}
		}
		username, _ := c.Get("username")
		uname, _ := username.(string)

		full := path
		if c.Request.URL.RawQuery != "" {
			full = full + "?" + c.Request.URL.RawQuery
		}

		global.DB.Create(&models.OperationLog{
			UserID:    userID,
			Username:  uname,
			Method:    c.Request.Method,
			Path:      full,
			Ip:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Status:    c.Writer.Status(),
			Latency:   time.Since(start).Milliseconds(),
		})
	}
}
