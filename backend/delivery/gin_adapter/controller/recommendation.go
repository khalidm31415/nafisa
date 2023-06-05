package controller

import (
	"backend/entity"
	"backend/internal_constant"
	"backend/usecase"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IRecommendation interface {
	MatchProfiles(c *gin.Context)
	View(c *gin.Context)
	Like(c *gin.Context)
	Pass(c *gin.Context)
}

type Recommendation struct {
	usecases usecase.Usecases
}

func NewRecommendation(usecases usecase.Usecases) IRecommendation {
	return &Recommendation{
		usecases: usecases,
	}
}

func (r *Recommendation) MatchProfiles(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := r.usecases.Recommendation.FindMatchingProfiles(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (r *Recommendation) View(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	profile, err := r.usecases.Recommendation.View(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (r *Recommendation) Like(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := r.usecases.Recommendation.Like(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (r *Recommendation) Pass(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := r.usecases.Recommendation.Pass(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
