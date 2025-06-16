package internal

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controller *Controller) func(*gin.Engine) {
	return func(r *gin.Engine) {
		routeGroup := r.Group("/brave")
		routeGroup.POST("/search", controller.Search)
	}
}
