package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func (m *middleware) CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     m.cfg.Cors.AllowOrigins,
		AllowMethods:     m.cfg.Cors.AllowMethods,
		AllowHeaders:     m.cfg.Cors.AllowHeaders,
		AllowCredentials: m.cfg.Cors.AllowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
