package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/module"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	cfg := configs.LoadConfig()

	db, err := configs.SetupDatabase(cfg.DB)
	if err != nil {
		fmt.Errorf("db error: %w", err)
	}

	gin.SetMode(cfg.Mode.Debug)
	r := gin.Default()
	api := r.Group("/api")

	authApp := module.BootstrapInit(db, cfg)
	module.AuthRouteRegister(api, authApp)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Server gagal dijalankan...")
	}
}
