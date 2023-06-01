package controller

import (
	"context"
	"nafisah/delivery/dto"
	"nafisah/entity"
	"nafisah/internal_constant"
	"nafisah/internal_error"
	"nafisah/usecase"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type IVerification interface {
	Verify(c *gin.Context)
}

type Verification struct {
	authMiddleware *jwt.GinJWTMiddleware
	identityKey    string
	usecases       usecase.Usecases
}

func NewVerification(usecases usecase.Usecases, identityKey string) IVerification {
	authMiddleware := setupAuthMiddleware(usecases, identityKey)
	return &Verification{
		authMiddleware: authMiddleware,
		identityKey:    identityKey,
		usecases:       usecases,
	}
}

func (v *Verification) Verify(c *gin.Context) {
	var input dto.VerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, _ := c.Get(v.identityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	err := v.usecases.Verification.Verify(ctx, input.UserID)
	if err == internal_error.ErrUnauthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
