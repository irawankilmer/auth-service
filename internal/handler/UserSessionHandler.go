package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
)

type UserSessionHandler struct {
	usService service.UserSessionService
}

func NewUserSessionHandler(usS service.UserSessionService) *UserSessionHandler {
	return &UserSessionHandler{usService: usS}
}

func (h *UserSessionHandler) RefreshToken(c *gin.Context) {
	res := response.NewResponder(c)
	var deviceID, ipAddress, userAgent string
	deviceID = "coba dulu"
	ipAddress = c.ClientIP()
	userAgent = c.Request.UserAgent()

	// cek refresh token
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		res.Unauthorized("refresh token tidak ditemukan")
		return
	}

	// refresh token
	token, err := h.usService.Refresh(c.Request.Context(), refreshToken, deviceID, ipAddress, userAgent)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.SetCookie("access_token", token.AccessToken, 900, "/", "localhost", true, true)           // 15 menit
	c.SetCookie("refresh_token", token.RefreshToken, 60*60*24*30, "/", "localhost", true, true) // 30 hari
	res.OK(token, "refresh token berhasil", nil)
}
