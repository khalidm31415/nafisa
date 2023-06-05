package controller

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/usecase"
	"context"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type IAuth interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	CurrentUserProfile(c *gin.Context)
	NewAdmin(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc
	AdminAuthMiddleware() gin.HandlerFunc
}

type Auth struct {
	authMiddleware           *jwt.GinJWTMiddleware
	adminTokenAuthMiddleware gin.HandlerFunc
	usecases                 usecase.Usecases
}

func NewAuth(usecases usecase.Usecases) IAuth {
	authMiddleware := setupAuthMiddleware(usecases)
	adminTokenAuthMiddleware := setupAdminTokenAuthMiddleware()
	return &Auth{
		authMiddleware:           authMiddleware,
		adminTokenAuthMiddleware: adminTokenAuthMiddleware,
		usecases:                 usecases,
	}
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return a.authMiddleware.MiddlewareFunc()
}

func (a *Auth) AdminAuthMiddleware() gin.HandlerFunc {
	return a.adminTokenAuthMiddleware
}

func (a *Auth) Signup(c *gin.Context) {
	var input dto.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.usecases.User.Signup(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (a *Auth) Login(c *gin.Context) {
	a.authMiddleware.LoginHandler(c)
}

func (a *Auth) Logout(c *gin.Context) {
	a.authMiddleware.LogoutHandler(c)
}

func (a *Auth) CurrentUserProfile(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	profile, err := a.usecases.User.CurrentUserProfile(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, profile)
}

func (a *Auth) NewAdmin(c *gin.Context) {
	var input dto.NewAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.usecases.Admin.NewAdmin(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
