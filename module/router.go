package module

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/handler"
)

func AuthRouteRegister(rg *gin.RouterGroup, app *BootstrapApp) {
	v := valigo.NewValigo()

	userHandler := handler.NewUserHandler(app.UserService, v)

	user := rg.Group("/users")
	user.GET("/", userHandler.GetAll)
	user.POST("/", userHandler.Create)
	user.GET("/:id", userHandler.FindByID)
	user.PATCH("/:id/username", userHandler.UsernameUpdate)
	user.PATCH("/:id/email", userHandler.EmailUpdate)
	user.DELETE("/:id", userHandler.Delete)
}
