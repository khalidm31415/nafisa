package controller

import (
	"nafisah/usecase"

	"github.com/gin-gonic/gin"
)

type IRecommendation interface {
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

func (r *Recommendation) View(c *gin.Context) {

}

func (r *Recommendation) Like(c *gin.Context) {

}

func (r *Recommendation) Pass(c *gin.Context) {

}
