package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/repository"
)

type Middleware interface {
	CORSMiddleware() gin.HandlerFunc
	AuthMiddleware() gin.HandlerFunc
	RoleMiddleware(matchType RoleMatchType, requiredRoles ...string) gin.HandlerFunc
	EmailVerifyMiddleware() gin.HandlerFunc
}

type middleware struct {
	cfg      *configs.AppConfig
	userRepo repository.UserRepository
}

func NewMiddleware(config *configs.AppConfig, u repository.UserRepository) Middleware {
	return &middleware{cfg: config, userRepo: u}
}
