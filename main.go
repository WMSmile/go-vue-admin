package main

import (
	"log"
	"strconv"

	"go-vue-admin/internal/casbin"
	"go-vue-admin/internal/config"
	"go-vue-admin/internal/db"
	"go-vue-admin/internal/global"
	"go-vue-admin/internal/initialize"
	"go-vue-admin/internal/router"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	global.Config = cfg

	if err := db.InitDB(); err != nil {
		log.Fatalf("init db: %v", err)
	}
	if err := casbin.InitCasbin(); err != nil {
		log.Fatalf("init casbin: %v", err)
	}
	initialize.Seed()

	r := router.SetupRouter()
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Println("server starting on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
