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
	UserService service.UserService
	Middleware  middleware.Middleware
}

func BootstrapInit(db *sql.DB, cfg *configs.AppConfig) *BootstrapApp {
	utilities := utils.NewUtility(cfg)

	usernameRepo := repository.NewUsernameHistoryRepository(db)
	emailRepo := repository.NewEmailHistoryRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo, roleRepo, usernameRepo, emailRepo, utilities, cfg)

	middlewares := middleware.NewMiddleware(cfg)
	return &BootstrapApp{
		Middleware:  middlewares,
		UserService: userService,
	}
}
