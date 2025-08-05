package brave

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/lokman928/search-proxy/internal/common/ratelimiter"
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module/brave/internal/controller"
	"github.com/lokman928/search-proxy/internal/module/brave/internal/router"
	"github.com/lokman928/search-proxy/internal/module/brave/internal/service"
)

type BraveModule struct {
	Service    *service.Service
	controller *controller.Controller
}

func NewBraveModule(lc fx.Lifecycle, cfg *config.Config) *BraveModule {
	braveCfg := cfg.Brave

	serviceCfg := &service.ServiceConfig{
		BaseUrl:     braveCfg.BaseUrl,
		ApiKey:      braveCfg.ApiKey,
		RateLimiter: ratelimiter.NewTokenRateLimiter(&braveCfg.RateLimit),
	}
	service := service.NewService(serviceCfg)

	controller := controller.NewController(service)

	return &BraveModule{
		Service:    service,
		controller: controller,
	}
}

func (b *BraveModule) RegisterRoutes() func(*gin.Engine) {
	return func(r *gin.Engine) {
		router.RegisterRoutes(b.controller)(r)
	}
}
