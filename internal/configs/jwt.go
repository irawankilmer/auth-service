package configs

import (
	"os"
	"time"
)

type JWTConfig struct {
	Secret         string
	AccessTokenTTL time.Duration
}

func getSecretOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
