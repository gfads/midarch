package env

import (
	"adaptive/adaptiveV3/environment/plugins/manager"
	"adaptive/adaptiveV3/sharedadaptive"
	"plugin"
	"shared"
	"sync"
)

var LocalPlugins map[string]plugin.Plugin
var RemotePlugins map[string]string

type Environment struct {}

func (e Environment) StartLocal(toManaging chan map[string]plugin.Plugin, wg *sync.WaitGroup) {
	for {
		x := e.senseLocal()
		toManaging <- x
	}
}

func (e Environment) StartRemote(toManaging chan map[string]string, wg *sync.WaitGroup) {
	for {
		x := e.senseRemote()
		toManaging <- x
	}
}

func (Environment) senseLocal() map[string]plugin.Plugin { // Executable
	newPlugins := map[string]plugin.Plugin{}

	// Plugin repository has not been created yet, then create it
	if LocalPlugins == nil {
		LocalPlugins = make(map[string]plugin.Plugin)
	}

	// Current number of plugins
	lenOld := len(LocalPlugins)
	oldPlugins := map[string]plugin.Plugin{}

	for i := range LocalPlugins {
		oldPlugins[i] = LocalPlugins[i]
	}

	// Load plugins from file
	LocalPlugins = manager.MyPlugin{}.LoadExecutables(sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)

	if len(LocalPlugins) == lenOld { // no new plugin
		return map[string]plugin.Plugin{}
	} else {
		for i := range LocalPlugins { // identify new plugins
			_, ok := oldPlugins[i]
			if !ok {
				newPlugins[i] = LocalPlugins[i]
			}
		}
	}

	// Check if the plugin has a symbol 'Behaviour' inside
	for i := range newPlugins {
		p := newPlugins[i]
		_, err := p.Lookup("Behaviour")
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), "'Behaviour' not found in plugin")
		}
	}

	return newPlugins
}

func (Environment) senseRemote() map[string]string { // Executable
	newPlugins := map[string]string{}

	// Plugin repository has not been created yet, then create it
	if RemotePlugins == nil {
		RemotePlugins = make(map[string]string)
	}

	// Current number of plugins
	lenOld := len(RemotePlugins)
	oldPlugins := map[string]string{}

	for i := range RemotePlugins {
		oldPlugins[i] = RemotePlugins[i]
	}

	// Load plugins from file
	RemotePlugins = manager.MyPlugin{}.LoadSources(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE)

	if len(RemotePlugins) == lenOld { // no new plugin
		return map[string]string{}
	} else {
		for i := range RemotePlugins { // identify new plugins
			_, ok := oldPlugins[i]
			if !ok {
				newPlugins[i] = RemotePlugins[i]
			}
		}
	}

	return newPlugins
}
