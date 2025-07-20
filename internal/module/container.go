package module

import (
	"github.com/lokman928/search-proxy/internal/module/brave"
	"github.com/lokman928/search-proxy/internal/module/healthcheck"
	"go.uber.org/fx"
)

type ModuleContainer struct {
	fx.In
	HealthcheckModule *healthcheck.HealthcheckModule
	BraveModule       *brave.BraveModule
}

func (m *ModuleContainer) GetModules() []IModule {
	return []IModule{
		m.HealthcheckModule,
		m.BraveModule,
	}
}
