package service

import (
	"context"
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

type UserService interface {
	GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error)
	Create(ctx context.Context, req request.UserCreateRequest) error
	FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error)
	EmailUpdate(ctx context.Context, user *response.UserDetailResponse, newEmail string) (bool, error)
	RolesUpdate(ctx context.Context, user *response.UserDetailResponse, newRoles []string) (bool, error)
	Delete(ctx context.Context, user *response.UserDetailResponse) error
}

type userService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	usernameRepo repository.UsernameHistoryRepository
	emailRepo    repository.EmailHistoryRepository
	utilities    utils.Utility
	config       *configs.AppConfig
	evService    EmailVerificationService
}

func NewUserService(
	ur repository.UserRepository, rp repository.RoleRepository, un repository.UsernameHistoryRepository,
	er repository.EmailHistoryRepository, ut utils.Utility, cfg *configs.AppConfig, ev EmailVerificationService,
) UserService {
	return &userService{
		userRepo: ur, roleRepo: rp, usernameRepo: un, emailRepo: er, utilities: ut, config: cfg, evService: ev,
	}
}

func (s *userService) GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error) {
	return s.userRepo.GetAll(ctx, limit, offset)
}

func (s *userService) Create(ctx context.Context, req request.UserCreateRequest) error {
	// cek roles
	roles, err := s.roleRepo.CheckRoles(ctx, req.Roles)
	if err != nil {
		return err
	}

	// cek email dari table users
	emailExists, err := s.userRepo.CheckEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return apperror.New(apperror.CodeEmailConflict, "email tidak dapat digunakan", err)
	}

	// cek email dari table email_history
	emailHistoryExists, err := s.emailRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailHistoryExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah tidak dapat digunakan", err)
	}

	// generate token version
	tokenVersion, err := s.utilities.UUIDGenerate()
	if err != nil {
		return apperror.New("[UUID_GENERATED_VALIED]", "gagal generate UUID", err, http.StatusInternalServerError)
	}

	// create user
	userID := s.utilities.ULIDGenerate()
	user := model.UserModel{
		ID:             userID,
		Username:       nil,
		Email:          req.Email,
		Password:       nil,
		TokenVersion:   tokenVersion,
		EmailVerified:  false,
		CreatedByAdmin: true,
		GoogleID:       nil,
		Profile: model.ProfileModel{
			ID:       s.utilities.ULIDGenerate(),
			UserID:   userID,
			FullName: &req.FullName,
			Address:  nil,
			Gender:   nil,
			Image:    nil,
		},
		Roles: roles,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return err
	}

	// kirim verifikasi email, dan atur waktu kadaluarsa token selama 7 hari
	if _, err := s.evService.SendVerification(ctx, &user, "verify-register-by-admin", "register", 168*time.Hour); err != nil {
		return err
	}

	return nil
}

func (s *userService) FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *userService) EmailUpdate(ctx context.Context, user *response.UserDetailResponse, newEmail string) (bool, error) {
	if user.Email == newEmail {
		return false, nil
	}

	// cek email dari tabel users
	emailChange, err := s.userRepo.EmailChange(ctx, user, newEmail)
	if err != nil {
		return false, err
	}
	if emailChange {
		return false, apperror.New(apperror.CodeEmailConflict, "email sudah terdaftar", err)
	}

	// cek email dari tabel email_history
	emailHistory, err := s.emailRepo.IsEmailExists(ctx, newEmail)
	if err != nil {
		return false, err
	}
	if emailHistory {
		return false, apperror.New(apperror.CodeEmailConflict, "email sudah terdaftar", err)
	}

	// update email
	if err := s.userRepo.EmailUpdate(ctx, user, newEmail); err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) RolesUpdate(ctx context.Context, user *response.UserDetailResponse, newRoles []string) (bool, error) {
	// cek role baru apakah tersedia di database
	newRolesCheck, err := s.roleRepo.CheckRoles(ctx, newRoles)
	if err != nil {
		return false, err
	}

	// bandingkan role lama dan baru
	if s.roleRepo.RoleIDsEqual(user.Roles, newRolesCheck) {
		return false, nil
	}

	// update roles
	if err := s.userRepo.RoleUpdate(ctx, user, newRolesCheck); err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) Delete(ctx context.Context, user *response.UserDetailResponse) error {
	return s.userRepo.Delete(ctx, user)
}
