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
	USService   service.UserSessionService
	CFG         *configs.AppConfig
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
	usRepo := repository.NewUserSessionRepository(db)

	evService := service.NewEmailVerificationService(evRepo, mail, utilities, cfg.Mail, userRepo, usernameRepo)
	userService := service.NewUserService(userRepo, roleRepo, usernameRepo, emailRepo, utilities, cfg, evService)
	authService := service.NewAuthService(authRepo, utilities, cfg, userRepo, roleRepo, usernameRepo, emailRepo, evService, usRepo)
	usService := service.NewUserSessionService(usRepo, utilities, cfg)

	middlewares := middleware.NewMiddleware(cfg, userRepo)
	return &BootstrapApp{
		AuthService: authService,
		Middleware:  middlewares,
		UserService: userService,
		EVService:   evService,
		USService:   usService,
	}
}
