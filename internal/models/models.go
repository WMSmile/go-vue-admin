package models

import "time"

// User <-> Role (many-to-many via user_roles)
// Role <-> Menu (many-to-many via role_menus)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Email     string    `gorm:"size:128" json:"email"`
	Phone     string    `gorm:"size:32" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int       `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Roles     []Role    `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Keyword     string    `gorm:"size:64;uniqueIndex;not null" json:"keyword"` // casbin sub
	Description string    `gorm:"size:255" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	Menus       []Menu    `gorm:"many2many:role_menus;" json:"menus"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Menu types: 1 = catalog, 2 = menu(page), 3 = button(permission)
type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ParentID   uint      `gorm:"default:0;index" json:"parentId"`
	Name       string    `gorm:"size:64" json:"name"`       // route name
	Title      string    `gorm:"size:64;not null" json:"title"` // display name
	Icon       string    `gorm:"size:64" json:"icon"`
	Path       string    `gorm:"size:255" json:"path"`       // frontend route path
	Component  string    `gorm:"size:255" json:"component"`  // vue component
	Sort       int       `gorm:"default:0" json:"sort"`
	Type       int       `gorm:"default:2" json:"type"`
	Permission string    `gorm:"size:64" json:"permission"` // e.g. user:add
	Api        string    `gorm:"size:255" json:"api"`        // backend api path
	Method     string    `gorm:"size:16" json:"method"`
	Status     int       `gorm:"default:1" json:"status"`
	Hidden     int       `gorm:"default:0" json:"hidden"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Children   []Menu    `gorm:"-" json:"children,omitempty"`
}

const (
	MenuTypeCatalog = 1
	MenuTypeMenu    = 2
	MenuTypeButton  = 3
)

// ApiKey is a machine credential that lets external callers (e.g. a skill/agent)
// access the API with the same RBAC permissions as the owning user. The raw Key
// is only returned at creation time and is never serialized in listings.
type ApiKey struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"userId"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name       string    `gorm:"size:64" json:"name"`
	Key        string    `gorm:"size:80;uniqueIndex;not null" json:"-"`
	Status     int       `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Scope      string    `gorm:"size:20;default:'all'" json:"scope"` // all | readonly
	ExpiresAt  *time.Time `json:"expiresAt"` // nil = never expire
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// OperationLog records every authenticated API request for audit purposes.
type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:64;index" json:"username"`
	Method    string    `gorm:"size:16" json:"method"`
	Path      string    `gorm:"size:255;index" json:"path"`
	Ip        string    `gorm:"size:64" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"userAgent"`
	Status    int       `json:"status"`
	Latency   int64     `json:"latency"` // milliseconds
	CreatedAt time.Time `json:"createdAt"`
}

// DictType groups a set of dictionary entries (e.g. gender, order_status).
type DictType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Status    int       `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DictData is a single key/value entry belonging to a DictType.
type DictData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TypeID    uint      `gorm:"index" json:"typeId"`
	TypeCode  string    `gorm:"size:64;index" json:"typeCode"`
	Label     string    `gorm:"size:128" json:"label"`
	Value     string    `gorm:"size:128" json:"value"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Task is a scheduled job (cron). Type 1 = HTTP request, Type 2 = built-in func.
type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:128;not null" json:"name"`
	Type      int        `gorm:"default:1" json:"type"` // 1 http, 2 func
	Spec      string     `gorm:"size:64" json:"spec"`   // 5-field cron
	Url       string     `gorm:"size:255" json:"url"`
	Method    string     `gorm:"size:16" json:"method"`
	Handler   string     `gorm:"size:64" json:"handler"`
	Status    int        `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Remark    string     `gorm:"size:255" json:"remark"`
	LastRunAt *time.Time `json:"lastRunAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// TaskLog records the result of a single task execution.
type TaskLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"index" json:"taskId"`
	TaskName   string    `gorm:"size:128" json:"taskName"`
	Status     int       `gorm:"default:1" json:"status"` // 1 success, 0 fail
	Output     string    `gorm:"type:text" json:"output"`
	DurationMs int64     `json:"durationMs"`
	CreatedAt  time.Time `json:"createdAt"`
}

// SysConfig is a key/value system parameter (e.g. site name, copyright).
// `config_key` is used as the column name to avoid the SQL reserved word `key`.
type SysConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"column:config_key;size:64;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Name      string    `gorm:"size:128;not null" json:"name"` // display / description
	Group     string    `gorm:"column:config_group;size:64" json:"group"` // grouping for the UI
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"` // 1 enabled, 0 disabled
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
