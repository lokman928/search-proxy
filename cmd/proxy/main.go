package main

import (
	"github.com/lokman928/search-proxy/internal/common/logger"
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module/brave"
	"github.com/lokman928/search-proxy/internal/module/healthcheck"
	"github.com/lokman928/search-proxy/internal/proxy"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.NewLogger("Fx").GetZapLogger()}
		}),
		fx.Provide(
			proxy.NewProxy,
			config.NewConfig,
			healthcheck.NewHealthcheckModule,
			brave.NewBraveModule,
		),
		fx.Invoke(func(p *proxy.Proxy) {}),
	)

	app.Run()
}
