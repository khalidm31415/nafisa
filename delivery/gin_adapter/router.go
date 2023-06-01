package gin_adapter

import (
	"nafisah/delivery/gin_adapter/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func SetupRouter(controllers *controller.Controllers) *gin.Engine {

	r := gin.Default()

	r.GET("/ping", Ping)

	r.POST("/auth/new-admin", controllers.Auth.AdminAuthMiddleware(), controllers.Auth.NewAdmin)
	r.POST("/auth/signup", controllers.Auth.Signup)
	r.POST("/auth/login", controllers.Auth.Login)
	r.POST("/auth/logout", controllers.Auth.Logout)
	r.GET("/auth/current-user", controllers.Auth.AuthMiddleware(), controllers.Auth.CurrentUser)

	r.POST("/verification/verify", controllers.Auth.AuthMiddleware(), controllers.Verification.Verify)

	r.GET("/recommendation", controllers.Auth.AuthMiddleware(), controllers.Recommendation.View)
	r.GET("/recommendation/like", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Like)
	r.GET("/recommendation/pass", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Pass)

	return r
}
