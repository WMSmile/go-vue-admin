package router

import (
	"net/http"
	"time"

	"go-vue-admin/internal/controller"
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if global.Config != nil && global.Config.Server.Mode != "" {
		gin.SetMode(global.Config.Server.Mode)
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	// public
	api.POST("/auth/login", controller.Login)

	// protected (API-Key or JWT + Casbin double check)
	auth := api.Group("")
	auth.Use(middleware.APIKeyAuth(), middleware.CasbinAuth(), middleware.OperationLog())
	{
		auth.GET("/menus/tree", controller.MenuTree)
		auth.GET("/auth/me", controller.Me)
		auth.POST("/auth/logout", controller.Logout)
		auth.GET("/dashboard", controller.Dashboard)

		auth.GET("/users", controller.ListUsers)
		auth.POST("/users", controller.CreateUser)
		auth.PUT("/users/:id", controller.UpdateUser)
		auth.DELETE("/users/:id", controller.DeleteUser)

		auth.GET("/roles", controller.ListRoles)
		auth.POST("/roles", controller.CreateRole)
		auth.PUT("/roles/:id", controller.UpdateRole)
		auth.DELETE("/roles/:id", controller.DeleteRole)
		auth.POST("/roles/menus", controller.AssignMenus)

		auth.GET("/menus", controller.ListMenus)

		auth.GET("/apikeys", controller.ListApiKeys)
		auth.POST("/apikeys", controller.CreateApiKey)
		auth.DELETE("/apikeys/:id", controller.DeleteApiKey)

		auth.GET("/operation-logs", controller.ListOperationLogs)
		auth.DELETE("/operation-logs/:id", controller.DeleteOperationLog)
		auth.DELETE("/operation-logs", controller.ClearOperationLogs)

		auth.GET("/dict-types", controller.ListDictTypes)
		auth.POST("/dict-types", controller.CreateDictType)
		auth.PUT("/dict-types/:id", controller.UpdateDictType)
		auth.DELETE("/dict-types/:id", controller.DeleteDictType)

		auth.GET("/dict-data", controller.ListDictData)
		auth.POST("/dict-data", controller.CreateDictData)
		auth.PUT("/dict-data/:id", controller.UpdateDictData)
		auth.DELETE("/dict-data/:id", controller.DeleteDictData)

		auth.GET("/tasks", controller.ListTasks)
		auth.POST("/tasks", controller.CreateTask)
		auth.PUT("/tasks/:id", controller.UpdateTask)
		auth.DELETE("/tasks/:id", controller.DeleteTask)
		auth.POST("/tasks/:id/toggle", controller.ToggleTask)
		auth.POST("/tasks/:id/run", controller.RunTask)
		auth.GET("/task-logs", controller.ListTaskLogs)
		auth.POST("/menus", controller.CreateMenu)
		auth.PUT("/menus/:id", controller.UpdateMenu)
		auth.DELETE("/menus/:id", controller.DeleteMenu)
	}
	return r
}
