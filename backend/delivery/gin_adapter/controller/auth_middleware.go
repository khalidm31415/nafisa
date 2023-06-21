package controller

import (
	"backend/internal_constant"
	"backend/usecase"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func setupAuthMiddleware(usecases usecase.Usecases) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("BACKEND_JWT_SECRET_KEY")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["userId"].(string)

			user, err := usecases.User.FindByUserID(c.Request.Context(), userID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			c.Set(internal_constant.GinIdentityKey, user)
			c.Next()
		} else {
			// Handle token validation error
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not valid"})
				return
			}
		}
	}
}

func setupAdminTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Validate the token
		if token != fmt.Sprintf("Bearer %s", os.Getenv("BACKEND_ADMIN_TOKEN")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Admin token"})
			c.Abort()
			return
		}

		// Proceed to the next middleware or route handler
		c.Next()
	}
}
