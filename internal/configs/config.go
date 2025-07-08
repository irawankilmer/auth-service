package configs

import (
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	DB       DBConfig
	Mode     GinModeConfig
	Server   ServerPortConfig
	Password Password
	Cors     CORSConfig
	JWT      JWTConfig
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
		Mode: GinModeConfig{
			Debug: getModeOrDefault("GIN_MODE", "debug"),
		},
		Server: ServerPortConfig{
			Port: getPortOrDefault("APP_PORT", "8080"),
		},
		Password: Password{
			Default: getPasswordOrDefault("PASSWORD_DEFAULT", "123456Aa*"),
		},
		Cors: CORSConfig{
			AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
			AllowMethods:     strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ","),
			AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ","),
			AllowCredentials: os.Getenv("CORS_ALLOW_CREDENTIALS") == "true",
		},
		JWT: JWTConfig{
			Secret:         getSecretOrDefault("JWT_SECRET", "default-secret"),
			AccessTokenTTL: 15 * time.Minute,
		},
	}
}
