package brave

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/common"
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
		RateLimiter: common.NewTokenRateLimiter(&braveCfg.RateLimit),
	}
	service := internal.NewService(serviceCfg)

	controller := internal.NewController(service)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Initializing Brave module...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Brave module...")
			return nil
		},
	})

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
