package env

import (
	manager2 "adaptive/adaptiveV2/environment/plugins/manager"
	"plugin"
	"shared"
	"sync"
)

var Plugins map[string]plugin.Plugin

type Environment struct{}

func (e Environment) Start(toManaging chan map[string]plugin.Plugin, wg *sync.WaitGroup) {
	for {
		x := e.sense()
		toManaging <- x
	}
}

func (Environment) sense() map[string]plugin.Plugin {
	newPlugins := map[string]plugin.Plugin{}

	// Plugin repository has not been created yet, then create it
	if Plugins == nil {
		Plugins = make(map[string]plugin.Plugin)
	}

	// Current number of plugins
	lenOld := len(Plugins)
	oldPlugins := map[string]plugin.Plugin{}

	for i := range Plugins {
		oldPlugins[i] = Plugins[i]
	}

	// Load plugins from file
	Plugins = manager2.MyPlugin{}.LoadPlugins()

	if len(Plugins) == lenOld { // no new plugin
		return map[string]plugin.Plugin{}
	} else {
		for i := range Plugins {  // identify new plugins
			_, ok := oldPlugins[i]
			if !ok {
				newPlugins[i] = Plugins[i]
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