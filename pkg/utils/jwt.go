package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/irawankilmer/auth-service/internal/configs"
	"time"
)

func (u *utility) JWTGenerate(userID, tokenVersion string, emailVerified bool, roles []string, cfg *configs.AppConfig) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":        userID,
		"token_version":  tokenVersion,
		"email_verified": emailVerified,
		"roles":          roles,
		"exp":            now.Add(cfg.JWT.AccessTokenTTL).Unix(),
		"iat":            now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := cfg.JWT.Secret

	return token.SignedString([]byte(secret))
}
