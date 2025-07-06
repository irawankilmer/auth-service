package main

import (
	"github.com/irawankilmer/auth-service/database"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	cfg := configs.LoadConfig()
	db, err := configs.SetupDatabase(cfg.DB)
	if err != nil {
		log.Fatal("koneksi ke database gagal:", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migrasi gagal: %v", err)
	}
}
