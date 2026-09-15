package casbin

import (
	"log"

	"go-vue-admin/internal/global"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// InitCasbin wires the gorm adapter to the enforcer and loads policies.
func InitCasbin() error {
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		return err
	}
	enf, err := casbin.NewEnforcer(global.Config.Casbin.Model, adapter)
	if err != nil {
		return err
	}
	if err := enf.LoadPolicy(); err != nil {
		return err
	}
	global.Enforcer = enf
	log.Println("casbin enforcer ready")
	return nil
}
