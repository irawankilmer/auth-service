package module

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/handler"
	"github.com/irawankilmer/auth-service/internal/middleware"
)

func AuthRouteRegister(rg *gin.RouterGroup, app *BootstrapApp) {
	v := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(app.AuthService, v)
	userHandler := handler.NewUserHandler(app.UserService, v)
	rg.Use(app.Middleware.CORSMiddleware())

	// role middleware
	saa := app.Middleware.RoleMiddleware(middleware.MatchAny, "super admin", "admin")

	// ===> auth routes
	auth := rg.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)

	// auth middleware
	auth.Use(app.Middleware.AuthMiddleware())
	auth.POST("/logout", authHandler.Logout)
	// ===> end auth routes

	// ===> users routes
	user := rg.Group("/users")
	user.Use(app.Middleware.AuthMiddleware())
	user.GET("", saa, userHandler.GetAll)
	user.POST("", saa, userHandler.Create)
	user.GET("/:id", saa, userHandler.FindByID)
	user.PATCH("/:id/username", saa, userHandler.UsernameUpdate)
	user.PATCH("/:id/email", saa, userHandler.EmailUpdate)
	user.PATCH("/:id/reset-password", saa, userHandler.PasswordReset)
	user.PATCH("/:id/roles-update", saa, userHandler.RoleUpdate)
	user.DELETE("/:id", saa, userHandler.Delete)
	// ===> end users routes
}
