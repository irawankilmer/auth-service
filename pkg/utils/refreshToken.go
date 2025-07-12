package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gogaruda/apperror"
)

func (u *utility) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (u *utility) RefreshTokenGenerate() (string, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal membuat refresh token", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
