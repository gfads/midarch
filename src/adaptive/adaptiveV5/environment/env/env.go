package env

import (
	"adaptive/adaptiveV5/environment/plugins/manager"
)

type EnvironmentInfo struct {
	Plugins map[string]string
}

type Environment interface {
	SensePlugins(string) EnvironmentInfo
}

type EnvironmentImpl struct {
	Info EnvironmentInfo
}

func NewEnvironment() Environment {
	var r Environment

	// Initialize folders
	plugins := make(map[string]string)

	info := EnvironmentInfo{Plugins: plugins}
	r = &EnvironmentImpl{Info: info}

	return r
}

func (e *EnvironmentImpl) SensePlugins(d string) EnvironmentInfo {

	// Load new plugins from files
	newPlugins := manager.MyPlugin{}.LoadSources(d)

	if len(newPlugins) != len(e.Info.Plugins) { // no new plugins
		for i := range newPlugins { // identify new plugins
			_, ok := e.Info.Plugins[i]
			if !ok {
				e.Info.Plugins[i] = newPlugins[i]
			}
		}
	}

	r := EnvironmentInfo{Plugins: e.Info.Plugins}

	return r
}
