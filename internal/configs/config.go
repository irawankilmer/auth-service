package configs

import "os"

type AppConfig struct {
	DB DBConfig
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		DB: DBConfig{
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
		},
	}
}
