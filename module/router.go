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

	r.Use(app.Middleware.CORSMiddleware())

	// role middleware
	saa := app.Middleware.RoleMiddleware(middleware.MatchAny, "super admin", "admin")

	// ===> auth routes
	auth := r.Group("/api/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/verify-email", emailVerifyHandler.VerifyEmail)
	auth.POST("/verify-register", emailVerifyHandler.VerifyRegister)
	auth.POST("/verify-register-resend", emailVerifyHandler.VerifyRegisterResend)

	// auth middleware
	auth.Use(app.Middleware.AuthMiddleware())
	auth.GET("/me", authHandler.Me)
	auth.POST("/logout", authHandler.Logout)
	// ===> end auth routes

	// ===> users routes
	user := r.Group("/users")
	user.Use(app.Middleware.AuthMiddleware())
	user.GET("", saa, userHandler.GetAll)
	user.POST("", saa, userHandler.Create)
	user.GET("/:id", saa, userHandler.FindByID)
	user.PATCH("/:id/email", saa, userHandler.EmailUpdate)
	user.PATCH("/:id/roles-update", saa, userHandler.RoleUpdate)
	user.DELETE("/:id", saa, userHandler.Delete)
	// ===> end users routes
}
