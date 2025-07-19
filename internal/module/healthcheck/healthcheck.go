package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/module/healthcheck/internal"
	"go.uber.org/fx"
)

type HealthcheckModule struct {
	controller *internal.Controller
}

func NewHealthcheckModule(lc fx.Lifecycle) *HealthcheckModule {
	controller := internal.NewController()

	return &HealthcheckModule{
		controller: controller,
	}
}

func (b *HealthcheckModule) RegisterRoutes() func(*gin.Engine) {
	return func(r *gin.Engine) {
		internal.RegisterRoutes(b.controller)(r)
	}
}
