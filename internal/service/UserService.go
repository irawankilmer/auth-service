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

func (s *userService) Delete(ctx context.Context, user *response.UserDetailResponse) error {
	return s.userRepo.Delete(ctx, user)
}
