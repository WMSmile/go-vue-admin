package db

import (
	"log"

	"go-vue-admin/internal/global"
	"go-vue-admin/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB opens the database and runs auto migration.
func InitDB() error {
	cfg := global.Config.Database

	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	default:
		dialector = sqlite.Open(cfg.DSN)
	}

	var err error
	global.DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}

	if err := global.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Menu{}, &models.ApiKey{}, &models.OperationLog{}, &models.DictType{}, &models.DictData{}); err != nil {
		return err
	}
	log.Printf("database connected (%s)", cfg.Driver)
	return nil
}
