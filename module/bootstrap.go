package module

import (
	"database/sql"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/middleware"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/utils"
)

type BootstrapApp struct {
	AuthService service.AuthService
	Middleware  middleware.Middleware
	UserService service.UserService
}

func BootstrapInit(db *sql.DB, cfg *configs.AppConfig) *BootstrapApp {
	utilities := utils.NewUtility(cfg)

	authRepository := repository.NewAuthRepository(db)
	usernameRepo := repository.NewUsernameHistoryRepository(db)
	emailRepo := repository.NewEmailHistoryRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(authRepository, utilities, cfg, userRepo)
	userService := service.NewUserService(userRepo, roleRepo, usernameRepo, emailRepo, utilities, cfg)

	middlewares := middleware.NewMiddleware(cfg, userRepo)
	return &BootstrapApp{
		AuthService: authService,
		Middleware:  middlewares,
		UserService: userService,
	}
}
