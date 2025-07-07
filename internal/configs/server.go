package configs

import "os"

type ServerPortConfig struct {
	Port string
}

func getPortOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
