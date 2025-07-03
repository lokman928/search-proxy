package brave

import (
	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/common/ratelimiter"
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module/brave/internal"
	"go.uber.org/fx"
)

type BraveModule struct {
	Service    *internal.Service
	controller *internal.Controller
}

func NewBraveModule(lc fx.Lifecycle, cfg *config.Config) *BraveModule {
	braveCfg := cfg.Brave

	serviceCfg := &internal.ServiceConfig{
		BaseUrl:     braveCfg.BaseUrl,
		ApiKey:      braveCfg.ApiKey,
		RateLimiter: ratelimiter.NewTokenRateLimiter(&braveCfg.RateLimit),
	}
	service := internal.NewService(serviceCfg)

	controller := internal.NewController(service)

	return &BraveModule{
		Service:    service,
		controller: controller,
	}
}

func (b *BraveModule) RegisterRoutes() func(*gin.Engine) {
	return func(r *gin.Engine) {
		internal.RegisterRoutes(b.controller)(r)
	}
}
