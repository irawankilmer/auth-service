package configs

import "os"

type AppConfig struct {
	DB       DBConfig
	Mode     GinModeConfig
	Server   ServerPortConfig
	Password Password
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
	}
}
