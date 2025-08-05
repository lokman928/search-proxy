package router

import (
	"github.com/gin-gonic/gin"

	"github.com/lokman928/search-proxy/internal/module/brave/internal/controller"
)

func RegisterRoutes(controller *controller.Controller) func(*gin.Engine) {
	return func(r *gin.Engine) {
		routeGroup := r.Group("/brave")
		routeGroup.POST("/search", controller.Search)
	}
}
