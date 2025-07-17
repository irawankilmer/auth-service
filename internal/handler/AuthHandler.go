package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
	validates   *valigo.Valigo
	userService service.UserService
	cfg         *configs.AppConfig
}

func NewAuthHandler(as service.AuthService, v *valigo.Valigo, u service.UserService, cfg *configs.AppConfig) *AuthHandler {
	return &AuthHandler{authService: as, validates: v, userService: u, cfg: cfg}
}

// Me godoc
// @Summary Ambil data user login
// @Description Mengambil data pengguna berdasarkan token
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	res := response.NewResponder(c)

	// cek user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		res.Unauthorized("user_id tidak ada di context")
		return
	}

	// ambil data user/me
	user, err := h.authService.Me(c.Request.Context(), userID.(string))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(user, "query ok", nil)
}

// Login godoc
// @Summary Login user
// @Description Login user dan generate token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.LoginRequest

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}

	// login
	token, err := h.authService.Login(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("access_token", token.AccessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", token.RefreshToken, 7*24*3600, "/", "", false, true)
	res.OK(token, "login berhasil", nil)
}

// Logout godoc
// @Summary Logout dari 1 device
// @Description Menghapus access & refresh token dari 1 device
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	res := response.NewResponder(c)

	// cek refresh token
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		res.Unauthorized("refresh token tidak ditemukan di cookie")
		return
	}

	// logout
	err = h.authService.Logout(c.Request.Context(), refreshToken)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// hapus cookie
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	res.OK(nil, "Logout berhasil", nil)
}

// LogoutAll godoc
// @Summary Logout dari semua device
// @Description Logout semua sesi user (semua device)
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/auth/logout-all-devices [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	res := response.NewResponder(c)

	// ambil user_id dari middleware JWT
	userID, exists := c.Get("user_id")
	if !exists {
		res.Unauthorized("user_id tidak ditemukan di context")
		return
	}

	// logout all devices
	if err := h.authService.LogoutAllDevices(c.Request.Context(), userID.(string)); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// hapus cookie
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	res.OK(nil, "Logout dari semua device berhasil", nil)
}

// Register godoc
// @Summary Registrasi user baru
// @Description Mendaftarkan user baru dan mengirim token verifikasi
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Data registrasi user baru"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.RegisterRequest
	req.Roles = []string{"tamu"}

	// validasi
	if !h.validates.ValigoJSON(c, &req) {
		return
	}
	errMap := make(map[string]string)
	if req.Password != req.ConfirmPassword {
		errMap["confirm_password"] = "konfirmasi password tidak cocok"
	}
	if !h.validates.ValigoBusiness(c, &req, errMap) {
		return
	}

	// registrasi
	token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("verify_email", token, 1800, "/", "", false, true)
	res.OK(token, "registrasi berhasil", nil)
}
