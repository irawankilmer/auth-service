package configs

import "os"

type GinModeConfig struct {
	Debug string
}

func getModeOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
