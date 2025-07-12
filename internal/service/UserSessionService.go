package service

import (
	"context"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/model"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"net/http"
	"time"
)

type UserSessionService interface {
	Refresh(ctx context.Context, refreshToken, deviceID, ipAddress, userAgent string) (*response.LoginResponse, error)
}

type userSessionServiceImpl struct {
	usRepo    repository.UserSessionRepository
	utilities utils.Utility
	cfg       *configs.AppConfig
}

func NewUserSessionService(usR repository.UserSessionRepository, util utils.Utility, cfg *configs.AppConfig) UserSessionService {
	return &userSessionServiceImpl{usRepo: usR, utilities: util, cfg: cfg}
}

func (s *userSessionServiceImpl) Refresh(ctx context.Context, refreshToken, deviceID, ipAddress, userAgent string) (*response.LoginResponse, error) {
	// cek refresh token hash
	session, err := s.usRepo.FindRefreshToken(ctx, s.utilities.HashToken(refreshToken))
	if err != nil {
		return nil, err
	}

	// cek apakah refresh token revoked atau sudah kadaluarsa
	if session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return nil, apperror.New("[REFRESH_TOKEN_INVALID]", "refresh token invalid", nil, http.StatusUnauthorized)
	}

	// cek token version
	user, err := s.usRepo.GetTokenVersionByUserID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// revoked token lama
	if err := s.usRepo.Revoked(ctx, session.ID); err != nil {
		return nil, err
	}

	// ambil roles
	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	// generate JWT
	accessToken, err := s.utilities.JWTGenerate(user.ID, user.TokenVersion, user.EmailVerified, roles, s.cfg)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "gagal generate JWT", err)
	}

	// generate new refresh token
	newRefreshToken, err := s.utilities.RefreshTokenGenerate()
	if err != nil {
		return nil, err
	}

	// create refresh token
	if err := s.usRepo.Create(ctx, &model.UserSession{
		ID:               s.utilities.ULIDGenerate(),
		UserID:           user.ID,
		RefreshTokenHash: s.utilities.HashToken(newRefreshToken),
		DeviceID:         deviceID,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
	}); err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
