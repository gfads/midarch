package main

import (
	"adaptive/adaptiveV3/selfadaptivesystem/managed"
	"adaptive/adaptiveV3/selfadaptivesystem/managing"
	"plugin"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	managedSystem := managed.Managed{}
	managingSystem := managing.ManagingSystem{}

	chn01,chn02,chn03 := initialiseChannels()

	wg.Add(2)

	go managingSystem.StartLocal(chn01, chn02, chn03, &wg)
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