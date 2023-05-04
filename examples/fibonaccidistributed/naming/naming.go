package main

import (
	"fmt"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
	"sync"
)

func main() {
	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["srh"] = messages.EndPoint{Host: "0.0.0.0", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("naming.madl", args)

	fmt.Printf("Naming server is running at Port: %v \n", shared.NAMING_PORT)
	//evolutive.EvolutiveInjector{}.Start("srhtcp", 40*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
