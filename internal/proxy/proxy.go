package proxy

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module"
	"go.uber.org/fx"
)

type Proxy struct {
	config config.ProxyConfig
	module *module.ModuleContainer
}

func NewProxy(lc fx.Lifecycle, cfg *config.Config, module module.ModuleContainer) *Proxy {
	p := &Proxy{
		config: cfg.Server,
		module: &module,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting Proxy server...")
			p.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Proxy server...")
			return nil
		},
	})

	return p
}

func (p *Proxy) Start() {
	r := gin.Default()

	for _, mod := range p.module.GetModules() {
		mod.RegisterRoutes()(r)
	}

	addr := fmt.Sprintf(":%d", p.config.Port)
	go r.Run(addr)
}
