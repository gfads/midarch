package injector

import (
	"adaptive/adaptiveV3/environment/plugins/manager"
	"adaptive/adaptiveV3/sharedadaptive"
	"strconv"
	"sync"
	"time"
)

type PluginInjector struct{}

func (PluginInjector) Initialize() {
	manager := manager.MyPlugin{}
	manager.InitialiseRepository()
}

func (p PluginInjector) StartLocal(wg *sync.WaitGroup) {

	manager := manager.MyPlugin{}

	for i := 0; i < 20; i++ {
		pluginName := "behaviour" + strconv.Itoa(i)
		code := manager.GenerateSource(pluginName)
		manager.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, pluginName, code)
		manager.GenerateExecutable(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL, pluginName)

		time.Sleep(sharedadaptive.PLUGIN_GENERATION_TIME * time.Second) // Generate a new plugin each 20 seconds
	}
	wg.Done()
}

func (p PluginInjector) StartRemote(wg *sync.WaitGroup) {

	manager := manager.MyPlugin{}

	i := 0

	for {
		pluginName := "behaviour" + strconv.Itoa(i)
		code := manager.GenerateSource(pluginName)
		manager.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE, pluginName, code)

		time.Sleep(sharedadaptive.PLUGIN_GENERATION_TIME * time.Second) // Generate a new plugin each 20 seconds
		i++
	}
	wg.Done()
}
