package injector

import (
	"adaptive/adaptiveV2/environment/plugins/manager"
	"adaptive/adaptiveV2/sharedadaptive"
	"strconv"
	"sync"
	"time"
)

type PluginInjector struct{}

func (PluginInjector) Initialize() {
	manager := manager.MyPlugin{}
	manager.InitialiseRepository()
}

func (p PluginInjector) Start(wg *sync.WaitGroup) {

	manager := manager.MyPlugin{}

	for i := 0; i < 20; i++ {
		pluginName := "behaviour"+strconv.Itoa(i)
		code := manager.GenerateSource(pluginName)
		manager.SaveCode(pluginName, code)
		manager.GenerateExecutable(pluginName)

		time.Sleep(sharedadaptive.PLUGIN_GENERATION_TIME * time.Second) // Generate a new plugin each 20 seconds
	}
	wg.Done()
}

