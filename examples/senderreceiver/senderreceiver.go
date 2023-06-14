package main

import (
	"sync"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	evolutive "github.com/gfads/midarch/pkg/injector"
)

func main() {
	fe := frontend.NewFrontend()

	args := make(map[string]messages.EndPoint)

	fe.Deploy(frontend.DeployOptions{FileName: "senderreceiver.madl", Args: args})

	evolutive.EvolutiveInjector{}.Start("sender", "sender", 20*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
