package services

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/apps"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/bindings"
)

func (s *ServiceInstance) BindWithParams(app *apps.App, params string) *bindings.Binding {
	return bindings.Bind(s.Name, app.Name, params)
}

func (s *ServiceInstance) Bind(app *apps.App) *bindings.Binding {
	return bindings.Bind(s.Name, app.Name, "")
}
