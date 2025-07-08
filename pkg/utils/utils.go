package utils

import "github.com/irawankilmer/auth-service/internal/configs"

type Utility interface {
	ULIDGenerate() string
	HashGenerate(password string) (string, error)
	HashCompare(hash, password string) bool
	UUIDGenerate() (string, error)
	JWTGenerate(userID, tokenVersion string, isVerified bool, roles []string, cfg *configs.AppConfig) (string, error)
}

type utility struct {
	config *configs.AppConfig
}

func NewUtility(cfg *configs.AppConfig) Utility {
	return &utility{config: cfg}
}
