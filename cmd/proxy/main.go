package main

import (
	"github.com/lokman928/search-proxy/internal/config"
	"github.com/lokman928/search-proxy/internal/module/brave"
	"github.com/lokman928/search-proxy/internal/proxy"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			proxy.NewProxy,
			config.NewConfig,
			brave.NewBraveModule,
		),
		fx.Invoke(func(p *proxy.Proxy) {}),
	)

	app.Run()
}
