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
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (string, error)
	Logout(ctx context.Context, userID string) error
	Register(ctx context.Context, req request.RegisterRequest) error
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
}

func NewAuthService(ar repository.AuthRepository, ut utils.Utility, cfg *configs.AppConfig,
	ur repository.UserRepository, rp repository.RoleRepository,
	username repository.UsernameHistoryRepository, email repository.EmailHistoryRepository,
	ev EmailVerificationService,
) AuthService {
	return &authService{
		authRepo: ar, utility: ut, cfg: cfg, userRepo: ur, roleRepo: rp,
		usernameRepo: username, emailRepo: email, evService: ev,
	}
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	// Cek identifikasi
	user, err := s.authRepo.IdentifierCheck(ctx, req.Identifier)
	if err != nil {
		return "", err
	}

	// cek verifikasi email
	if !user.EmailVerified {
		return "", apperror.New("[EMAIL_NOT_VERIFY]", "email belum di verifikasi", err, http.StatusUnauthorized)
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

func (s *authService) Logout(ctx context.Context, userID string) error {
	// cek user
	user, err := s.authRepo.Me(ctx, userID)
	if err != nil {
		return err
	}

	// buat token_version baru
	newTokenVersion, err := s.utility.UUIDGenerate()
	if err != nil {
		return err
	}

	// update token_version di database
	if err := s.authRepo.UpdateTokenVersion(ctx, user.ID, newTokenVersion); err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest) error {
	// cek roles
	roles, err := s.roleRepo.CheckRoles(ctx, req.Roles)
	if err != nil {
		return err
	}

	// cek username dari table users
	usernameExists, err := s.userRepo.CheckUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return apperror.New(apperror.CodeUsernameConflict, "username tidak dapat digunakan", err)
	}

	// cek email dari table users
	emailExists, err := s.userRepo.CheckEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return apperror.New(apperror.CodeEmailConflict, "email tidak dapat digunakan", err)
	}

	// cek username history
	usernameHistoryExists, err := s.usernameRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return err
	}
	if usernameHistoryExists {
		return apperror.New(apperror.CodeUsernameConflict, "username sudah tidak dapat digunakan", err)
	}

	// cek email history
	emailHistoryExists, err := s.emailRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailHistoryExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah tidak dapat digunakan", err)
	}

	// generate password
	passHash, err := s.utility.HashGenerate(req.Password)
	if err != nil {
		return apperror.New("[CODE_GENERATE_HASH_INVALID]", "gagal generate hash", err, 505)
	}

	userID := s.utility.ULIDGenerate()
	user := model.UserModel{
		ID:             userID,
		Username:       &req.Username,
		Email:          req.Email,
		Password:       &passHash,
		TokenVersion:   nil,
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
		return err
	}

	// kirim verifikasi email
	if err := s.evService.SendVerification(ctx, user, "verify-email"); err != nil {
		return err
	}

	return nil
}

func (s *authService) Me(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	return s.authRepo.Me(ctx, userID)
}
