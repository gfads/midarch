package injector

import (
	"adaptive/adaptiveV6/environment/plugins/manager"
	"adaptive/adaptiveV6/sharedadaptive"
	"strconv"
	"sync"
	"time"
)

type PluginInjector struct{}

func (PluginInjector) Initialize() {
	manager := manager.MyPlugin{}
	manager.InitialiseRepository()
}

func (p PluginInjector) Start(t int, wg *sync.WaitGroup) {
	manager := manager.MyPlugin{}

	switch t {
	case sharedadaptive.LOCAL: // Generate Executable plugins
		i := 0
		for {
			plugin := manager.GenerateSource("behaviour" + strconv.Itoa(i))
			manager.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, plugin)
			manager.GenerateExecutable(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL, plugin)

			time.Sleep(sharedadaptive.PLUGIN_GENERATION_TIME * time.Second) // Generate a new plugin each XX seconds
			i++
		}
	case sharedadaptive.REMOTE: // Generate Source plugins
		i := 0
		for {
			p := manager.GenerateSource("behaviour" + strconv.Itoa(i))
			manager.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE, p)

			time.Sleep(sharedadaptive.PLUGIN_GENERATION_TIME * time.Second) // Generate a new plugin each XX seconds
			i++
		}
	}
	wg.Done()
}
