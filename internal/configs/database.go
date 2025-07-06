package configs

import (
	"database/sql"
	"fmt"
	"time"
)

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

func SetupDatabase(cfg DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("koneksi database gagal: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping ke database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
