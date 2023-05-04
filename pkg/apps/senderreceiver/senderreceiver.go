package main

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	evolutive "github.com/gfads/midarch/pkg/injector"
	"sync"
	"time"
)

func main() {
	fe := frontend.NewFrontend()

	args := make(map[string]messages.EndPoint)

	fe.Deploy("senderreceiver.madl", args)

	evolutive.EvolutiveInjector{}.Start("sender", 20*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
