package main

import (
	"fmt"
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	"shared"
	"sync"
	"time"
)

func main() {
	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: "localhost", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("calculatordistributedclientmid.madl", args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux, ok := namingProxy.Lookup("Calculator")
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Service 'Calculator' not found in Naming Service")
	}

	calc := aux.(*calculatorproxy.Calculatorproxy)
	for x := 0; x < 1000; x++ {
		fmt.Println("Result:", calc.Add(x, 1))
		time.Sleep(200 * time.Millisecond)
	}

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
