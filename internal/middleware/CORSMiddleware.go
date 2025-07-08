package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (m *middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowAll := len(m.cfg.Cors.AllowOrigins) == 1 && m.cfg.Cors.AllowOrigins[0] == "*"

		// Handle allowed origin
		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && contains(m.cfg.Cors.AllowOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		} else {
			// Origin tidak diizinkan, lanjut tapi tidak set CORS headers
			c.Next()
			return
		}

		// Common headers
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(m.cfg.Cors.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(m.cfg.Cors.AllowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", "600") // Optional: cache preflight 10 menit

		if m.cfg.Cors.AllowCredentials && !allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Fungsi pembantu agar rapi
func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
