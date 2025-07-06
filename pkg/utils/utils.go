package utils

import "github.com/irawankilmer/auth-service/internal/configs"

type Utility interface {
	ULIDGenerate() string
	HashGenerate(password string) (string, error)
	HashCompare(hash, password string) bool
}

type utility struct {
	config *configs.AppConfig
}

func NewUtility(cfg *configs.AppConfig) Utility {
	return &utility{config: cfg}
}
