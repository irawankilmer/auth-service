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

		// Ambil token dari Cookie atau Authorization Header
		var tokenStr string

		// 1. Coba dari Cookie
		if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
			tokenStr = cookieToken
		} else {
			// 2. Coba dari Authorization header
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

		// Ambil user_id
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			res.Unauthorized("token tidak memiliki user_id yang valid")
			return
		}

		// Ambil token_version
		tokenVersion, ok := claims["token_version"].(string)
		if !ok || tokenVersion == "" {
			res.Unauthorized("token_version tidak valid")
			return
		}

		// Ambil email_verified
		emailVerified, ok := claims["email_verified"].(bool)
		if !ok {
			res.Unauthorized("field email_verified tidak valid")
			return
		}

		// Ambil roles
		rolesClaim, ok := claims["roles"].([]interface{})
		if !ok {
			res.Unauthorized("format roles dalam token tidak valid")
			return
		}

		roles := make([]string, 0, len(rolesClaim))
		for _, r := range rolesClaim {
			roleStr, ok := r.(string)
			if !ok {
				res.Unauthorized("setiap role harus berupa string")
				return
			}
			roles = append(roles, roleStr)
		}

		// Cek token_version ke DB
		user, err := m.userRepo.FindUserByTokenVersion(c.Request.Context(), userID)
		if err != nil {
			res.Unauthorized("user tidak ditemukan")
			return
		}
		if user.TokenVersion == nil || *user.TokenVersion != tokenVersion {
			res.Unauthorized("token sudah tidak berlaku, silakan login ulang")
			return
		}

		// Simpan ke context
		c.Set("user_id", userID)
		c.Set("email_verified", emailVerified)
		c.Set("roles", roles)

		c.Next()
	}
}
