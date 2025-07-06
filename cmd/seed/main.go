package main

import (
	"fmt"
	"github.com/irawankilmer/auth-service/database/seeders"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/pkg/utils"
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

	u := utils.NewUtility(cfg)
	if err := seeders.SeedsRun(db, u); err != nil {
		log.Fatalf("seeding gagal:%v", err)
	}

	fmt.Println("seeder berhasil...")
}
