package module

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/handler"
	"github.com/irawankilmer/auth-service/internal/middleware"
)

func AuthRouteRegister(r *gin.Engine, app *BootstrapApp) {
	v := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(app.AuthService, v, app.UserService)
	userHandler := handler.NewUserHandler(app.UserService, v)
	emailVerifyHandler := handler.NewEmailVerificationHandler(app.EVService, v)
	uSessionHandler := handler.NewUserSessionHandler(app.USService)

	r.Use(app.Middleware.CORSMiddleware())

	// role middleware
	saa := app.Middleware.RoleMiddleware(middleware.MatchAny, "super admin", "admin")

	// ===> auth routes
	auth := r.Group("/api/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/logout-all-devices", authHandler.LogoutAll)
	auth.POST("/register", authHandler.Register)
	auth.POST("/verify-email", emailVerifyHandler.VerifyEmail)
	auth.POST("/verify-register-resend", emailVerifyHandler.VerifyRegisterResend)
	auth.POST("/verify-register-by-admin", emailVerifyHandler.VerifyRegisterByAdmin)
	auth.POST("/verify-register-by-admin-resend", emailVerifyHandler.VerifyRegisterByAdminResend)

	// auth middleware
	auth.Use(app.Middleware.AuthMiddleware())
	auth.GET("/me", authHandler.Me)
	// ===> end auth routes

	// refresh token
	refresh := r.Group("/api/refresh-token")
	refresh.POST("", uSessionHandler.RefreshToken)

	// ===> users routes
	user := r.Group("/api/users")
	user.Use(app.Middleware.AuthMiddleware())
	user.GET("", saa, userHandler.GetAll)
	user.POST("", saa, userHandler.Create)
	user.GET("/:id", saa, userHandler.FindByID)
	user.PATCH("/:id/email", saa, userHandler.EmailUpdate)
	user.PATCH("/:id/roles-update", saa, userHandler.RoleUpdate)
	user.DELETE("/:id", saa, userHandler.Delete)
	// ===> end users routes
}
