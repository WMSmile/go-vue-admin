package controller

import (
	"strconv"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"
	"go-vue-admin/internal/task"
	"go-vue-admin/internal/utils"

	"github.com/gin-gonic/gin"
)

// ListTasks returns a paginated, filterable list of tasks.
func ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	name := c.Query("name")
	status := c.Query("status")

	db := global.DB.Model(&models.Task{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			db = db.Where("status = ?", s)
		}
	}

	var total int64
	db.Count(&total)
	var list []models.Task
	db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	utils.Ok(c, gin.H{"list": list, "total": total})
}

type TaskReq struct {
	Name    string `json:"name" binding:"required"`
	Type    int    `json:"type"` // 1 http, 2 func
	Spec    string `json:"spec" binding:"required"`
	Url     string `json:"url"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
	Status  int    `json:"status"`
	Remark  string `json:"remark"`
}

// CreateTask creates a task and (if enabled) registers it in the scheduler.
func CreateTask(c *gin.Context) {
	var req TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "invalid params")
		return
	}
	if req.Type == 1 && req.Url == "" {
		utils.Fail(c, "HTTP 任务需填写 URL")
		return
	}
	if req.Type == 2 && req.Handler == "" {
		utils.Fail(c, "函数任务需填写 Handler")
		return
	}
	t := models.Task{
		Name:    req.Name,
		Type:    req.Type,
		Spec:    req.Spec,
		Url:     req.Url,
		Method:  req.Method,
		Handler: req.Handler,
		Status:  req.Status,
		Remark:  req.Remark,
	}
	global.DB.Create(&t)
	task.AddTask(t)
	utils.Ok(c, t)
}

// UpdateTask updates a task and refreshes the live scheduler entry.
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var t models.Task
	if err := global.DB.First(&t, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	var req TaskReq
	c.ShouldBindJSON(&req)
	t.Name = req.Name
	t.Type = req.Type
	t.Spec = req.Spec
	t.Url = req.Url
	t.Method = req.Method
	t.Handler = req.Handler
	t.Status = req.Status
	t.Remark = req.Remark
	global.DB.Save(&t)
	task.AddTask(t)
	utils.Ok(c, t)
}

// DeleteTask removes a task and its logs, and stops it in the scheduler.
func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	global.DB.Where("task_id = ?", id).Delete(&models.TaskLog{})
	global.DB.Delete(&models.Task{}, id)
	task.RemoveTask(uint(id))
	utils.Ok(c, nil)
}

// ToggleTask enables/disables a task.
func ToggleTask(c *gin.Context) {
	id := c.Param("id")
	var t models.Task
	if err := global.DB.First(&t, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	t.Status = 1 - t.Status
	global.DB.Save(&t)
	task.SetEnabled(t)
	utils.Ok(c, t)
}

// RunTask triggers a one-off execution of the task.
func RunTask(c *gin.Context) {
	id := c.Param("id")
	var t models.Task
	if err := global.DB.First(&t, id).Error; err != nil {
		utils.Fail(c, "not found")
		return
	}
	task.RunNow(t)
	utils.Ok(c, nil)
}

// ListTaskLogs returns paginated execution logs, optionally for one task.
func ListTaskLogs(c *gin.Context) {
	taskID := c.Query("taskId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	db := global.DB.Model(&models.TaskLog{})
	if taskID != "" {
		db = db.Where("task_id = ?", taskID)
	}
	var total int64
	db.Count(&total)
	var list []models.TaskLog
	db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	utils.Ok(c, gin.H{"list": list, "total": total})
}
