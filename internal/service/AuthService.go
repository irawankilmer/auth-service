package service

import (
	"context"
	"errors"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"net/http"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (string, error)
}

type authService struct {
	authRepo repository.AuthRepository
	utility  utils.Utility
	cfg      *configs.AppConfig
}

func NewAuthService(ar repository.AuthRepository, ut utils.Utility, cfg *configs.AppConfig) AuthService {
	return &authService{authRepo: ar, utility: ut, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	// Cek identifikasi
	user, err := s.authRepo.IdentifierCheck(ctx, req.Identifier)
	if err != nil {
		return "", err
	}

	// cek password
	if !s.utility.HashCompare(*user.Password, req.Password) {
		return "", apperror.New("[PASSWORD_INVALID]", "password salah", errors.New("Password salah"), http.StatusUnauthorized)
	}

	// buat uuid
	newTokenVersion, err := s.utility.UUIDGenerate()
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "token version gagal dibuat", err)
	}

	// ambil roles user
	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	// Generate token
	token, err := s.utility.JWTGenerate(user.ID, newTokenVersion, user.EmailVerified, roles, s.cfg)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "Generate token gagal", err)
	}

	// Update token_version
	if err := s.authRepo.UpdateTokenVersion(ctx, user.ID, newTokenVersion); err != nil {
		return "", err
	}

	return token, nil
}
