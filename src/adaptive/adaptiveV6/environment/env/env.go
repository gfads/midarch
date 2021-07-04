package env

import (
	"adaptive/adaptiveV6/environment/plugins/manager"
)

type Environment interface {
	SensePlugins(string) []manager.MyPlugin
}

type EnvironmentImpl struct {
}

func NewEnvironment() Environment {
	var r Environment

	r = &EnvironmentImpl{}

	return r
}

func (e *EnvironmentImpl) SensePlugins(d string) []manager.MyPlugin {

	// Load new plugins from files
	r := manager.MyPlugin{}.LoadSources(d)

	return r
}
