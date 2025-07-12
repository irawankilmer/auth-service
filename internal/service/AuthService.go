package service

import (
	"context"
	"errors"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/model"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest, userAgent, ipAddress string) (*response.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAllDevices(ctx context.Context, userID string) error
	Register(ctx context.Context, req request.RegisterRequest) (string, error)
	Me(ctx context.Context, userID string) (*response.UserDetailResponse, error)
}

type authService struct {
	authRepo     repository.AuthRepository
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	utility      utils.Utility
	cfg          *configs.AppConfig
	usernameRepo repository.UsernameHistoryRepository
	emailRepo    repository.EmailHistoryRepository
	evService    EmailVerificationService
	usRepo       repository.UserSessionRepository
}

func NewAuthService(ar repository.AuthRepository, ut utils.Utility, cfg *configs.AppConfig,
	ur repository.UserRepository, rp repository.RoleRepository,
	username repository.UsernameHistoryRepository, email repository.EmailHistoryRepository,
	ev EmailVerificationService, usR repository.UserSessionRepository,
) AuthService {
	return &authService{
		authRepo: ar, utility: ut, cfg: cfg, userRepo: ur, roleRepo: rp,
		usernameRepo: username, emailRepo: email, evService: ev, usRepo: usR,
	}
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest, userAgent, ipAddress string) (*response.LoginResponse, error) {
	// Cek identifikasi
	user, err := s.authRepo.IdentifierCheck(ctx, req.Identifier)
	if err != nil {
		return nil, err
	}

	// cek verifikasi email
	if !user.EmailVerified {
		return nil, apperror.New("[EMAIL_NOT_VERIFY]", "email belum di verifikasi", err, http.StatusUnauthorized)
	}

	// cek password
	if !s.utility.HashCompare(*user.Password, req.Password) {
		return nil, apperror.New("[PASSWORD_INVALID]", "password salah", errors.New("Password salah"), http.StatusUnauthorized)
	}

	// ambil roles user
	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	// Generate token
	token, err := s.utility.JWTGenerate(user.ID, user.TokenVersion, user.EmailVerified, roles, s.cfg)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternalError, "Generate token gagal", err)
	}

	// generate refresh token
	refreshToken, err := s.utility.RefreshTokenGenerate()
	if err != nil {
		return nil, err
	}

	// insert refresh token
	if err := s.usRepo.Create(ctx, &model.UserSession{
		ID:               s.utility.ULIDGenerate(),
		UserID:           user.ID,
		RefreshTokenHash: s.utility.HashToken(refreshToken),
		DeviceID:         "coba saja",
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
	}); err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	hashed := s.utility.HashToken(refreshToken)

	session, err := s.usRepo.FindRefreshToken(ctx, hashed)
	if err != nil {
		// alasan return nil (dianggap berhasil logout), agar status token tidak bocor, begitu :)
		return nil
	}
	if session.Revoked {
		return nil
	}

	return s.usRepo.Revoked(ctx, session.ID)
}

func (s *authService) LogoutAllDevices(ctx context.Context, userID string) error {
	// generate token version baru
	newTokenVersion, err := s.utility.UUIDGenerate()
	if err != nil {
		return apperror.New(apperror.CodeInternalError, "generate new token version gagal", err)
	}

	// update token version
	if err := s.authRepo.UpdateTokenVersion(ctx, userID, newTokenVersion); err != nil {
		return err
	}

	//revoke semua session
	if err := s.usRepo.RevokeAllSessionByUserID(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest) (string, error) {
	// cek roles
	roles, err := s.roleRepo.CheckRoles(ctx, req.Roles)
	if err != nil {
		return "", err
	}

	// cek username dari table users
	usernameExists, err := s.userRepo.CheckUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}
	if usernameExists {
		return "", apperror.New(apperror.CodeUsernameConflict, "username tidak dapat digunakan", err)
	}

	// cek email dari table users
	emailExists, err := s.userRepo.CheckEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", apperror.New(apperror.CodeEmailConflict, "email tidak dapat digunakan", err)
	}

	// cek username history
	usernameHistoryExists, err := s.usernameRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return "", err
	}
	if usernameHistoryExists {
		return "", apperror.New(apperror.CodeUsernameConflict, "username sudah tidak dapat digunakan", err)
	}

	// cek email history
	emailHistoryExists, err := s.emailRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if emailHistoryExists {
		return "", apperror.New(apperror.CodeEmailConflict, "email sudah tidak dapat digunakan", err)
	}

	// generate password
	passHash, err := s.utility.HashGenerate(req.Password)
	if err != nil {
		return "", apperror.New("[CODE_GENERATE_HASH_INVALID]", "gagal generate hash", err, 505)
	}

	// generate token version
	tokenVersion, err := s.utility.UUIDGenerate()
	if err != nil {
		if err != nil {
			return "", apperror.New("[UUID Generate VAILED]", "gagal membuat UUID", err, http.StatusInternalServerError)
		}
	}

	userID := s.utility.ULIDGenerate()
	user := model.UserModel{
		ID:             userID,
		Username:       &req.Username,
		Email:          req.Email,
		Password:       &passHash,
		TokenVersion:   tokenVersion,
		EmailVerified:  false,
		CreatedByAdmin: false,
		GoogleID:       nil,
		Profile: model.ProfileModel{
			ID:       s.utility.ULIDGenerate(),
			UserID:   userID,
			FullName: &req.FullName,
			Address:  nil,
			Gender:   nil,
			Image:    nil,
		},
		Roles: roles,
	}

	// register
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return "", err
	}

	// kirim verifikasi email
	token, err := s.evService.SendVerification(ctx, &user, "verify-email", "register", 30*time.Minute)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Me(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	return s.authRepo.Me(ctx, userID)
}
