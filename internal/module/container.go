package module

import (
	"github.com/lokman928/search-proxy/internal/module/brave"
	"go.uber.org/fx"
)

type ModuleContainer struct {
	fx.In
	BraveModule *brave.BraveModule
}

func (m *ModuleContainer) GetModules() []IModule {
	return []IModule{
		m.BraveModule,
	}
}
