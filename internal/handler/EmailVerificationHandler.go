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

// VerifyEmail godoc
// @Summary Verifikasi token email
// @Description Memverifikasi token yang dikirim melalui email saat registrasi
// @Tags Verifikasi Email
// @Accept json
// @Produce json
// @Param request body request.VerifyRequest true "Payload token verifikasi"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/auth/verify-email [post]
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

// VerifyRegisterResend godoc
// @Summary Kirim ulang token verifikasi email saat register
// @Description Mengirim ulang token verifikasi dari cookie "verify_email"
// @Tags Verifikasi Email
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/auth/verify-register-resend [post]
func (h *EmailVerificationHandler) VerifyRegisterResend(c *gin.Context) {
	res := response.NewResponder(c)
	token, err := c.Cookie("verify_email")
	if err != nil {
		res.Unauthorized("Kesalahan saat mengambil token verifikasi")
		return
	}
	ctx := c.Request.Context()

	// cek token dan ambil data user
	user, err := h.evService.CheckToken(ctx, token)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// kirim verifikasi
	newToken, err := h.evService.SendVerification(ctx, user, "verify-email", "register", 30*time.Minute)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("verify_email", newToken, 1800, "/", "", false, true)
	res.OK(newToken, "verifikasi email sudah dirikim ulang", nil)
}

// VerifyRegisterByAdmin godoc
// @Summary Verifikasi pendaftaran oleh admin
// @Description Admin melakukan verifikasi terhadap user menggunakan token
// @Tags Verifikasi Email
// @Accept json
// @Produce json
// @Param request body request.VerifyRegisterByAdminRequest true "Data verifikasi oleh admin"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/auth/verify-register-by-admin [post]
func (h *EmailVerificationHandler) VerifyRegisterByAdmin(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.VerifyRegisterByAdminRequest
	ctx := c.Request.Context()

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

	// cek token
	ev, err := h.evService.FindByToken(ctx, req.Token)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// update data registrasi
	if err := h.evService.UpdateRegisterByAdmin(ctx, &req, ev); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(c.Query("email"), "verifikasi berhasil", nil)
}

// VerifyRegisterByAdminResend godoc
// @Summary Kirim ulang token verifikasi pendaftaran oleh admin
// @Description Mengirim ulang token verifikasi ke user dari token yang dikirim admin
// @Tags Verifikasi Email
// @Accept json
// @Produce json
// @Param request body request.VerifyRegisterByAdminResendRequest true "Payload token lama dari admin"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/auth/verify-register-by-admin-resend [post]
func (h *EmailVerificationHandler) VerifyRegisterByAdminResend(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.VerifyRegisterByAdminResendRequest
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
	newToken, err := h.evService.SendVerification(ctx, user, "verify-register-by-admin", "register", 30*time.Minute)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("verify_email", newToken, 1800, "/", "", false, true)
	res.OK(newToken, "verifikasi email sudah dirikim ulang", nil)
}
