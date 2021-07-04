package main

import (
	"adaptive/adaptiveV2/environment/env"
	"adaptive/adaptiveV2/environment/plugins/injector"
	"adaptive/adaptiveV2/selfadaptivesystem/managed"
	"adaptive/adaptiveV2/selfadaptivesystem/managing"
	"plugin"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	environment := env.Environment{}
	managedSystem := managed.Managed{}
	managingSystem := managing.ManagingSystem{}
	inj := injector.PluginInjector{}

	chn01,chn02,chn03 := initialiseChannels()

	wg.Add(4)

	inj.Initialize()
	go inj.Start(&wg)
	go environment.Start(chn01, &wg)
	go managingSystem.Start(chn01, chn02, chn03, &wg)
	go managedSystem.Start(chn02, chn03, &wg)

	//go AdaptationGoals()

	wg.Wait()
}


func initialiseChannels() (chan map[string]plugin.Plugin, chan int, chan func()){
	chn01 := make(chan map[string]plugin.Plugin)
	chn02 := make(chan int)
	chn03 := make(chan func())

	return chn01,chn02,chn03
}