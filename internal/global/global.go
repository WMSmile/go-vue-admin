package global

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-vue-admin/internal/config"
)

var (
	DB       *gorm.DB
	Enforcer *casbin.Enforcer
	Config   *config.Config
)
