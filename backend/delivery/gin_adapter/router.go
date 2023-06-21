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
	r.GET("/verification/unverified-users", controllers.Auth.AuthMiddleware(), controllers.Verification.GetUnverifiedUsers)
	r.POST("/verification/verify", controllers.Auth.AuthMiddleware(), controllers.Verification.Verify)

	r.POST("/auth/signup", controllers.Auth.Signup)
	r.POST("/auth/login", controllers.Auth.Login)
	r.GET("/auth/google/login", controllers.Auth.GoogleLogin)
	r.GET("/auth/google/callback", controllers.Auth.HandleGoogleCallback)
	r.POST("/auth/logout", controllers.Auth.Logout)

	r.POST("/profile", controllers.Auth.AuthMiddleware(), controllers.Auth.CompleteProfile)
	r.GET("/profile", controllers.Auth.AuthMiddleware(), controllers.Auth.CurrentUserProfile)

	r.POST("/recommendation/match-profile", controllers.Auth.AuthMiddleware(), controllers.Recommendation.MatchProfiles)
	r.GET("/recommendation", controllers.Auth.AuthMiddleware(), controllers.Recommendation.View)
	r.POST("/recommendation/like", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Like)
	r.POST("/recommendation/pass", controllers.Auth.AuthMiddleware(), controllers.Recommendation.Pass)

	return r
}
