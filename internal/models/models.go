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
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
