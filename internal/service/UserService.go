package service

import (
	"context"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/dto/response"
	"github.com/irawankilmer/auth-service/internal/model"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/pkg/utils"
)

type UserService interface {
	GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error)
	Create(ctx context.Context, req request.UserCreateRequest) error
	FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error)
	UsernameUpdate(ctx context.Context, user *response.UserDetailResponse, newUsername string) (bool, error)
	EmailUpdate(ctx context.Context, user *response.UserDetailResponse, newEmail string) (bool, error)
	Delete(ctx context.Context, user *response.UserDetailResponse) error
}

type userService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	usernameRepo repository.UsernameHistoryRepository
	emailRepo    repository.EmailHistoryRepository
	utilities    utils.Utility
}

func NewUserService(
	ur repository.UserRepository,
	rp repository.RoleRepository,
	un repository.UsernameHistoryRepository,
	er repository.EmailHistoryRepository,
	ut utils.Utility) UserService {
	return &userService{
		userRepo: ur, roleRepo: rp, usernameRepo: un, emailRepo: er, utilities: ut,
	}
}

func (s *userService) GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error) {
	return s.userRepo.GetAll(ctx, limit, offset)
}

func (s *userService) Create(ctx context.Context, req request.UserCreateRequest) error {
	roles, err := s.roleRepo.CheckRoles(ctx, req.Roles)
	if err != nil {
		return err
	}

	usernameExists, err := s.userRepo.CheckUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return apperror.New(apperror.CodeUsernameConflict, "username tidak dapat digunakan", err)
	}

	emailExists, err := s.userRepo.CheckEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return apperror.New(apperror.CodeEmailConflict, "email tidak dapat digunakan", err)
	}

	usernameHistoryExists, err := s.usernameRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return err
	}
	if usernameHistoryExists {
		return apperror.New(apperror.CodeUsernameConflict, "username sudah tidak dapat digunakan", err)
	}

	emailHistoryExists, err := s.emailRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if emailHistoryExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah tidak dapat digunakan", err)
	}

	passHash, err := s.utilities.HashGenerate(req.Password)
	if err != nil {
		return apperror.New("[CODE_GENERATE_HASH_INVALID]", "gagal generate hash", err, 505)
	}

	userID := s.utilities.ULIDGenerate()
	user := model.UserModel{
		ID:             userID,
		Username:       &req.Username,
		Email:          req.Email,
		Password:       &passHash,
		TokenVersion:   nil,
		EmailVerified:  false,
		CreatedByAdmin: true,
		GoogleID:       nil,
		Profile: model.ProfileModel{
			ID:       s.utilities.ULIDGenerate(),
			UserID:   userID,
			FullName: req.Profile.FullName,
			Address:  nil,
			Gender:   nil,
			Image:    nil,
		},
		Roles: roles,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (s *userService) FindByID(ctx context.Context, userID string) (*response.UserDetailResponse, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *userService) UsernameUpdate(ctx context.Context, user *response.UserDetailResponse, newUsername string) (bool, error) {
	// cek apakah username diubah?
	currentUsername := ""
	if user.Username != nil {
		currentUsername = *user.Username
	}

	if currentUsername == newUsername {
		return false, nil
	}

	// cek username dari table users
	usernameChange, err := s.userRepo.UsernameChange(ctx, user, newUsername)
	if err != nil {
		return false, err
	}

	if usernameChange {
		return false, apperror.New(apperror.CodeUsernameConflict, "username sudah terdaftar", err)
	}

	// cek username dari tabel username_history
	usernameExists, err := s.usernameRepo.IsUsernameExists(ctx, newUsername)
	if err != nil {
		return false, err
	}
	if usernameExists {
		return false, apperror.New(apperror.CodeUsernameConflict, "username sudah terdaftar", err)
	}

	// update username
	if err := s.userRepo.UsernameUpdate(ctx, user, newUsername); err != nil {
		return false, err
	}

	return true, nil
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

func (s *userService) Delete(ctx context.Context, user *response.UserDetailResponse) error {
	return s.userRepo.Delete(ctx, user)
}
