package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/irawankilmer/auth-service/docs"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/module"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Auth Service API
// @version 1.0
// @description Dokumentasi API untuk Auth Service
// @termsOfService http://yourapp.com/terms/

// @contact.name Kirdun Developer
// @contact.url https://yourapp.com
// @contact.email dev@yourapp.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <token>
func main() {
	_ = godotenv.Load()
	cfg := configs.LoadConfig()

	db, err := configs.SetupDatabase(cfg.DB)
	if err != nil {
		fmt.Errorf("db error: %w", err)
	}

	gin.SetMode(cfg.Mode.Debug)
	r := gin.Default()

	authApp := module.BootstrapInit(db, cfg)
	module.AuthRouteRegister(r, authApp)

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Server gagal dijalankan...")
	}
}
