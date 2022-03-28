package main

import (
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	"shared"
	"sync"
)

func main() {
	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["srh"] = messages.EndPoint{Host: "localhost", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("naming.madl", args)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
