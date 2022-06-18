package main

import (
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	evolutive "injector"
	"sync"
	"time"
)

func main() {
	fe := frontend.NewFrontend()

	args := make (map[string]messages.EndPoint)

	fe.Deploy("senderreceiver.madl", args)

	evolutive.EvolutiveInjector{}.Start("sender", 20*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
