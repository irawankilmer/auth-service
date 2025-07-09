package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/irawankilmer/auth-service/pkg/response"
)

func (m *middleware) EmailVerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.NewResponder(c)
		isVerifiedRaw, exists := c.Get("email_verified")
		if !exists {
			res.Forbidden("status verifikasi email tidak tersedia")
			return
		}

		emailVerified, ok := isVerifiedRaw.(bool)
		if !ok || !emailVerified {
			res.Forbidden("email belum diverifikasi")
			return
		}

		c.Next()
	}
}
