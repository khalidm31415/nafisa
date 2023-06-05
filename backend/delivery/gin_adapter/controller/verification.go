package controller

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/internal_error"
	"backend/usecase"
	"context"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type IVerification interface {
	GetUnverifiedUsers(c *gin.Context)
	Verify(c *gin.Context)
}

type Verification struct {
	authMiddleware *jwt.GinJWTMiddleware
	usecases       usecase.Usecases
}

func NewVerification(usecases usecase.Usecases) IVerification {
	authMiddleware := setupAuthMiddleware(usecases)
	return &Verification{
		authMiddleware: authMiddleware,
		usecases:       usecases,
	}
}

func (v *Verification) GetUnverifiedUsers(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)

	users, err := v.usecases.Verification.GetUnverifiedUsers(ctx)
	if err == internal_error.ErrUnauthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (v *Verification) Verify(c *gin.Context) {
	var input dto.VerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, _ := c.Get(internal_constant.GinIdentityKey)
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
	if err = v.usecases.Admin.IndexProfile(ctx, input.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
