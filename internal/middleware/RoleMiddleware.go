package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/irawankilmer/auth-service/pkg/response"
	"strings"
)

type RoleMatchType int

const (
	MatchAny RoleMatchType = iota
	MatchAll
)

func normalizeRoles(roles []string) []string {
	normalized := make([]string, 0, len(roles))
	for _, r := range roles {
		r = strings.ToLower(strings.TrimSpace(r))
		if r != "" {
			normalized = append(normalized, r)
		}
	}

	return normalized
}

func matchRoles(userRoles, requiredRoles []string, matchType RoleMatchType) bool {
	roleSet := make(map[string]bool, len(userRoles))
	for _, r := range userRoles {
		roleSet[r] = true
	}

	switch matchType {
	case MatchAny:
		for _, r := range requiredRoles {
			if roleSet[r] {
				return true
			}
		}
		return false

	case MatchAll:
		for _, r := range requiredRoles {
			if !roleSet[r] {
				return false
			}
		}
		return true

	default:
		return false
	}
}

func (m *middleware) RoleMiddleware(matchType RoleMatchType, requiredRoles ...string) gin.HandlerFunc {
	normalizedRequired := normalizeRoles(requiredRoles)

	return func(c *gin.Context) {
		res := response.NewResponder(c)

		rawRoles, exists := c.Get("roles")
		if !exists {
			res.Forbidden("akses ditolak: roles tidak ditemukan dalam token")
			return
		}

		userRoles, ok := rawRoles.([]string)
		if !ok {
			res.Forbidden("akses ditolak: format roles tidak valid")
			return
		}

		normalizedUserRoles := normalizeRoles(userRoles)

		if !matchRoles(normalizedUserRoles, normalizedRequired, matchType) {
			res.Forbidden("akses ditolak: role tidak memenuhi syarat")
			return
		}

		c.Next()
	}
}
