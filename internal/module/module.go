package module

import "github.com/gin-gonic/gin"

type IModule interface {
	RegisterRoutes() func(*gin.Engine)
}
