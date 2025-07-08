package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	dto "github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
	validates   *valigo.Valigo
}

func NewAuthHandler(as service.AuthService, v *valigo.Valigo) *AuthHandler {
	return &AuthHandler{authService: as, validates: v}
}

func (h *AuthHandler) Login(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.LoginRequest

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}

	// login
	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(dto.LoginResponse{
		Token: token,
	}, "login berhasil", nil)
}
