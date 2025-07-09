package module

import (
	"database/sql"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/middleware"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/mailer"
	"github.com/irawankilmer/auth-service/pkg/utils"
)

type BootstrapApp struct {
	AuthService service.AuthService
	Middleware  middleware.Middleware
	UserService service.UserService
	EVService   service.EmailVerificationService
}

func BootstrapInit(db *sql.DB, cfg *configs.AppConfig) *BootstrapApp {
	utilities := utils.NewUtility(cfg)

	mail := mailer.NewMailer(cfg.Mail)
	authRepo := repository.NewAuthRepository(db)
	usernameRepo := repository.NewUsernameHistoryRepository(db)
	emailRepo := repository.NewEmailHistoryRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	evRepo := repository.NewEmailVerificationRepository(db)

	userService := service.NewUserService(userRepo, roleRepo, usernameRepo, emailRepo, utilities, cfg)
	evService := service.NewEmailVerificationService(evRepo, mail, utilities, cfg.Mail, userService)
	authService := service.NewAuthService(authRepo, utilities, cfg, userRepo, roleRepo, usernameRepo, emailRepo, evService)

	middlewares := middleware.NewMiddleware(cfg, userRepo)
	return &BootstrapApp{
		AuthService: authService,
		Middleware:  middlewares,
		UserService: userService,
		EVService:   evService,
	}
}
