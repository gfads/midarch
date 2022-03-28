package main

import (
	"fmt"
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/generic"
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
	args["crh"] = messages.EndPoint{Host: "localhost", Port: shared.NAMING_PORT}
	args["srh"] = messages.EndPoint{Host: "localhost", Port: shared.CALCULATOR_PORT}

	// Deploy configuration
	fe.Deploy("calculatordistributedservermid.madl", args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	// Create proxy to calculatorimpl
	calcProxy := calculatorproxy.NewCalculatorProxy(generic.ProxyConfig{
		Host: shared.CALCULATOR_HOST,
		Port: shared.CALCULATOR_PORT,
	})

	// Register calculatorimpl in Lookup
	ok := namingProxy.Register("Calculator", calcProxy)

	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "'Calculator' already registered in the Naming Server")
	}

	fmt.Printf("Calculator server is running at Port: %v \n", shared.CALCULATOR_PORT)

	//evolutive.EvolutiveInjector{}.Start("sender", 40*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
