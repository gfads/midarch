package env

import (
	"adaptive/adaptiveV5/environment/plugins/manager"
	"adaptive/adaptiveV5/sharedadaptive"
	"plugin"
	"shared"
)

type EnvironmentInfo struct {
	SourcePlugins     map[string]string
	ExecutablePlugins map[string]plugin.Plugin
}

type Environment interface {
	Sense(int, string) EnvironmentInfo
}

type EnvironmentImpl struct {
	Info EnvironmentInfo
}

func NewEnvironment() Environment {
	var r Environment

	sourcePlugins := make(map[string]string)
	executablePlugins := make(map[string]plugin.Plugin)
	info := EnvironmentInfo{SourcePlugins: sourcePlugins, ExecutablePlugins: executablePlugins}

	r = &EnvironmentImpl{Info: info}

	return r
}

func (e *EnvironmentImpl) Sense(t int, d string) EnvironmentInfo {
	r := EnvironmentInfo{}

	switch t {
	case sharedadaptive.SOURCE:
		if e.Info.SourcePlugins == nil {
			e.Info.SourcePlugins = make(map[string]string)
		}

		// Load plugins from file
		newPlugins := manager.MyPlugin{}.LoadSources(d)

		if len(newPlugins) != len(e.Info.SourcePlugins) { // no new plugin
			for i := range newPlugins { // identify new plugins
				_, ok := e.Info.SourcePlugins[i]
				if !ok {
					e.Info.SourcePlugins[i] = newPlugins[i]
				}
			}
		}
		r = EnvironmentInfo{SourcePlugins: e.Info.SourcePlugins}
	case sharedadaptive.EXECUTABLE:
		if e.Info.ExecutablePlugins == nil {
			e.Info.ExecutablePlugins = make(map[string]plugin.Plugin)
		}

		newExecutablePlugins := manager.MyPlugin{}.LoadExecutables(d)

		if len(newExecutablePlugins) != len(e.Info.ExecutablePlugins) { // no new plugin
			for i := range newExecutablePlugins { // identify new plugins
				_, ok := e.Info.ExecutablePlugins[i]
				if !ok {
					e.Info.ExecutablePlugins[i] = newExecutablePlugins[i]
				}
			}
		}

		// Check if the plugins have a symbol 'Behaviour' inside
		for i := range e.Info.ExecutablePlugins {
			p := newExecutablePlugins[i]
			_, err := p.Lookup("Behaviour")
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "'Behaviour' not found in plugin")
			}
		}
		r = EnvironmentInfo{ExecutablePlugins: e.Info.ExecutablePlugins}
	}
	return r
}

