package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
)

type EmailVerificationHandler struct {
	evService service.EmailVerificationService
}

func NewEmailVerificationHandler(ev service.EmailVerificationService) *EmailVerificationHandler {
	return &EmailVerificationHandler{evService: ev}
}

func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	res := response.NewResponder(c)

	// periksa token
	token := c.Query("token")
	if token == "" {
		res.BadRequest(nil, "token tidak boleh kosong")
		return
	}

	// verifikasi token
	if err := h.evService.VerifyToken(c.Request.Context(), token); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(nil, "Email berhasil di verifikas", nil)
}
