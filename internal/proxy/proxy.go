package proxy

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/common/logger"
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module"
	"github.com/lokman928/search-proxy/internal/proxy/internal/middleware"
	"go.uber.org/fx"
)

type Proxy struct {
	config config.ProxyConfig
	module *module.ModuleContainer
	logger logger.ILogger
}

func NewProxy(lc fx.Lifecycle, cfg *config.Config, module module.ModuleContainer) *Proxy {
	p := &Proxy{
		config: cfg.Server,
		module: &module,
		logger: logger.NewLogger("proxy"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return p
}

func (p *Proxy) Start() {
	r := gin.New()

	r.Use(middleware.LogRequest(p.logger), gin.Recovery())

	for _, mod := range p.module.GetModules() {
		mod.RegisterRoutes()(r)
	}

	addr := fmt.Sprintf(":%d", p.config.Port)
	go r.Run(addr)
}
