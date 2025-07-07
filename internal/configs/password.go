package configs

import "os"

type Password struct {
	Default string
}

func getPasswordOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
