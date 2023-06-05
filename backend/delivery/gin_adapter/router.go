package gin_adapter

import (
	"backend/delivery/gin_adapter/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func SetupRouter(controllers *controller.Controllers) *gin.Engine {

	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	r.GET("/ping", Ping)

	r.POST("/auth/new-admin", controllers.Auth.AdminAuthMiddleware(), controllers.Auth.NewAdmin)
	r.POST("/auth/signup", controllers.Auth.Signup)
	r.POST("/auth/login", controllers.Auth.Login)
	r.POST("/auth/logout", controllers.Auth.Logout)
	r.GET("/auth/current-user", controllers.Auth.AuthMiddleware(), controllers.Auth.CurrentUser)

	r.GET("/verification/unverified-users", controllers.Auth.AuthMiddleware(), controllers.Verification.GetUnverifiedUsers)
	r.POST("/verification/verify", controllers.Auth.AuthMiddleware(), controllers.Verification.Verify)

	r.POST("/recommendation/match-profile", controllers.Auth.AuthMiddleware(), controllers.Recommendation.MatchProfiles)
	r.GET("/recommendation", controllers.Auth.AuthMiddleware(), controllers.Recommendation.View)
	r.POST("/recommendation/like", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Like)
	r.POST("/recommendation/pass", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Pass)

	return r
}
