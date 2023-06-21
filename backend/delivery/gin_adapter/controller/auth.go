package controller

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	"backend/usecase"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type IAuth interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	GoogleLogin(c *gin.Context)
	HandleGoogleCallback(c *gin.Context)
	CompleteProfile(c *gin.Context)
	Logout(c *gin.Context)
	CurrentUserProfile(c *gin.Context)
	NewAdmin(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc
	AdminAuthMiddleware() gin.HandlerFunc
}

type Auth struct {
	googleOauthConfig        oauth2.Config
	authMiddleware           gin.HandlerFunc
	adminTokenAuthMiddleware gin.HandlerFunc
	usecases                 usecase.Usecases
}

func NewAuth(googleOauthConfig oauth2.Config, usecases usecase.Usecases) IAuth {
	authMiddleware := setupAuthMiddleware(usecases)
	adminTokenAuthMiddleware := setupAdminTokenAuthMiddleware()
	return &Auth{
		googleOauthConfig:        googleOauthConfig,
		authMiddleware:           authMiddleware,
		adminTokenAuthMiddleware: adminTokenAuthMiddleware,
		usecases:                 usecases,
	}
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return a.authMiddleware
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
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.usecases.User.FindByUsername(c.Request.Context(), input.Username)
	if err != nil && err != internal_error.ErrUsernameNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err == internal_error.ErrUsernameNotFound {
		c.Status(http.StatusUnauthorized)
		return
	}

	if user.Password == nil {
		err := errors.New("password not set")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		err := errors.New("wrong password")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": user.ID,
		})
	secret := []byte(os.Getenv("BACKEND_JWT_SECRET_KEY"))
	tokenString, err := t.SignedString(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, false)
	c.Status(http.StatusOK)
}

func (a *Auth) GoogleLogin(c *gin.Context) {
	url := a.googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *Auth) HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := a.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := a.googleOauthConfig.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer res.Body.Close()
	// Parse and use the user information returned in the response
	var userInfo map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gmail := userInfo["email"].(string)
	user, err := a.usecases.User.FindByGmail(c.Request.Context(), gmail)
	if err != nil && err != internal_error.ErrGmailNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err == internal_error.ErrGmailNotFound {
		err := a.usecases.User.GoogleSignup(c.Request.Context(), gmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": user.ID,
		})
	secret := []byte(os.Getenv("BACKEND_JWT_SECRET_KEY"))
	tokenString, err := t.SignedString(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, false)

	if !user.Profile.IsProfileComplete {
		c.Redirect(http.StatusFound, "/profile")
		return
	}

	c.Redirect(http.StatusFound, "/recommendation")
}

func (a *Auth) CompleteProfile(c *gin.Context) {
	var input dto.CompleteProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := a.usecases.User.CompleteProfile(ctx, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (a *Auth) Logout(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
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
