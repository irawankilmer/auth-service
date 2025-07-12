package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/irawankilmer/auth-service/pkg/response"
)

func (m *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.NewResponder(c)

		// Ambil token dari Cookie atau Header
		var tokenStr string
		if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
			tokenStr = cookieToken
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				res.Unauthorized("token tidak ditemukan di cookie maupun header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				res.Unauthorized("format Authorization harus: Bearer {token}")
				return
			}
			tokenStr = parts[1]
		}

		// Parse token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			res.Unauthorized("token tidak valid atau sudah kadaluwarsa")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			res.Unauthorized("klaim token tidak valid")
			return
		}

		// Validasi exp
		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			res.Unauthorized("token sudah kadaluwarsa")
			return
		}

		// Ambil isi claim
		userID, _ := claims["user_id"].(string)
		tokenVersion, _ := claims["token_version"].(string)
		emailVerified, _ := claims["email_verified"].(bool)

		// Ambil roles (array of string)
		roles := []string{}
		if rawRoles, ok := claims["roles"].([]interface{}); ok {
			for _, r := range rawRoles {
				if str, ok := r.(string); ok {
					roles = append(roles, str)
				}
			}
		}

		if userID == "" || tokenVersion == "" {
			res.Unauthorized("token tidak memiliki user_id atau token_version yang valid")
			return
		}

		// Simpan data token ke context
		c.Set("user_id", userID)
		c.Set("token_version", tokenVersion)
		c.Set("email_verified", emailVerified)
		c.Set("roles", roles)

		c.Next()
	}
}
