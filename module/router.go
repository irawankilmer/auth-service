package module

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/handler"
)

func AuthRouteRegister(rg *gin.RouterGroup, app *BootstrapApp) {
	v := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(app.AuthService, v)
	userHandler := handler.NewUserHandler(app.UserService, v)
	rg.Use(app.Middleware.CORSMiddleware())

	auth := rg.Group("/auth")
	auth.POST("/login", authHandler.Login)

	user := rg.Group("/users")
	user.Use(app.Middleware.AuthMiddleware())
	user.GET("", userHandler.GetAll)
	user.POST("", userHandler.Create)
	user.GET("/:id", userHandler.FindByID)
	user.PATCH("/:id/username", userHandler.UsernameUpdate)
	user.PATCH("/:id/email", userHandler.EmailUpdate)
	user.PATCH("/:id/reset-password", userHandler.PasswordReset)
	user.PATCH("/:id/roles-update", userHandler.RoleUpdate)
	user.DELETE("/:id", userHandler.Delete)
}
