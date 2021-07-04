package main

import (
	"adaptive/adaptiveV3/environment/env"
	"adaptive/adaptiveV3/environment/plugins/injector"
	"adaptive/adaptiveV3/selfadaptivesystem/managing"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	environment := env.Environment{}
	managingSystem := managing.ManagingSystem{}
	inj := injector.PluginInjector{}

	chn01 := initialiseChannels()

	wg.Add(3)

	inj.Initialize()
	go inj.StartRemote(&wg)
	go environment.StartRemote(chn01, &wg)
	go managingSystem.StartRemote(chn01, &wg)

	//go AdaptationGoals()

	wg.Wait()
}


func initialiseChannels() chan map[string]string {
	chn01 := make(chan map[string]string)

	return chn01
}