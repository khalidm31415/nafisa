package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"nafisah/delivery/dto"
	"nafisah/entity"
	"nafisah/internal_error"
	"nafisah/usecase"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthMiddleware(usecases usecase.Usecases, identityKey string) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(os.Getenv("API_JWT_SECRET_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals dto.LoginInput
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			givenUsername := loginVals.Username
			givenPassword := loginVals.Password

			user, err := usecases.User.FindByUsername(context.Background(), givenUsername)

			if user == nil {
				return nil, internal_error.ErrInternalError
			}

			if errors.Is(err, internal_error.ErrUsernameNotFound) || user == nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(givenPassword)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user.ID, nil

		},
		PayloadFunc: func(userID interface{}) jwt.MapClaims {
			if userID != nil {
				return jwt.MapClaims{
					identityKey: userID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			userID := claims[identityKey].(string)
			user, _ := usecases.User.FindByUserID(c.Request.Context(), userID)
			return user
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*entity.User); ok {
				return true
			}
			return false
		},
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost",
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware
}

func setupAdminTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Validate the token
		if token != fmt.Sprintf("Bearer %s", os.Getenv("API_ADMIN_TOKEN")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Admin token"})
			c.Abort()
			return
		}

		// Proceed to the next middleware or route handler
		c.Next()
	}
}
