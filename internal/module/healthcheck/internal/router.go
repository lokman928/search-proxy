package internal

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controller *Controller) func(*gin.Engine) {
	return func(r *gin.Engine) {
		r.GET("/healthz", controller.HealthCheck)
	}
}
