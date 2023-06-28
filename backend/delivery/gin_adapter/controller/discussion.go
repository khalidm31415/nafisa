package controller

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/usecase"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IDiscussion interface {
	SendMessage(c *gin.Context)
	EndDiscussion(c *gin.Context)
}

type Discussion struct {
	usecases usecase.Usecases
}

func NewDiscussion(usecases usecase.Usecases) IDiscussion {
	return &Recommendation{
		usecases: usecases,
	}
}

func (r *Recommendation) SendMessage(c *gin.Context) {
	var input dto.SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := r.usecases.Discussion.SendMessage(ctx, input.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (r *Recommendation) EndDiscussion(c *gin.Context) {
	u, _ := c.Get(internal_constant.GinIdentityKey)
	user, _ := u.(*entity.User)
	ctx := context.WithValue(c.Request.Context(), internal_constant.ContextUserKey, user)
	if err := r.usecases.Discussion.EndDiscussion(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
