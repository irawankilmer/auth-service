package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/irawankilmer/auth-service/internal/configs"
)

type Middleware interface {
	CORSMiddleware() gin.HandlerFunc
}

type middleware struct {
	cfg *configs.AppConfig
}

func NewMiddleware(config *configs.AppConfig) Middleware {
	return &middleware{cfg: config}
}
