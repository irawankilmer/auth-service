package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
	"time"
)

type EmailVerificationHandler struct {
	evService service.EmailVerificationService
	validates *valigo.Valigo
}

func NewEmailVerificationHandler(ev service.EmailVerificationService, v *valigo.Valigo) *EmailVerificationHandler {
	return &EmailVerificationHandler{evService: ev, validates: v}
}

func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.VerifyRequest

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}

	// verifikasi token
	if err := h.evService.VerifyToken(c.Request.Context(), req.Token); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(nil, "Email berhasil di verifikas", nil)
}

func (h *EmailVerificationHandler) VerifyRegister(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.VerifyRegisterRequest

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}

	// validasi kecocokan password
	errMap := map[string]string{}
	if req.Password != req.PasswordConfirm {
		errMap["password_confirm"] = "konfirmasi password salah"
	}
	if !h.validates.ValigoBusiness(c, &req, errMap) {
		return
	}

	res.OK(c.Query("email"), "verifikasi berhasil", nil)
}

func (h *EmailVerificationHandler) VerifyRegisterResend(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.VerifyRequest
	ctx := c.Request.Context()

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}

	// cek token dan ambil data user
	user, err := h.evService.CheckToken(ctx, req.Token)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// kirim verifikasi
	token, err := h.evService.SendVerification(ctx, user, "verify-email", "register", 30*time.Minute)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("verify_email", token, 86400, "/", "", false, true)
	res.OK(token, "verifikasi email sudah dirikim ulang", nil)
}
