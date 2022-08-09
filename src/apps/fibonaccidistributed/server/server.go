package main

import (
	"fmt"
	"gmidarch/development/components/proxies/fibonacciProxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/generic"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	evolutive "injector"
	"shared"
	"sync"
	"time"
)

func main() {
	// Wait for namingserver to get up
	time.Sleep(10 * time.Second)

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: "namingserver", Port: shared.NAMING_PORT}
	args["srh"] = messages.EndPoint{Host: "0.0.0.0", Port: shared.CALCULATOR_PORT}

	// Deploy configuration
	fe.Deploy("FibonacciDistributedServerMid.madl", args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	// Create proxy to Fibonacci
	fibonacciProxy := fibonacciProxy.NewFibonacciProxy(generic.ProxyConfig{
		Host: shared.CALCULATOR_HOST,
		Port: shared.CALCULATOR_PORT,
	})

	// Register Fibonacci in Lookup
	ok := namingProxy.Register("Fibonacci", fibonacciProxy)

	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "'Fibonacci' already registered in the Naming Server")
	}

	fmt.Printf("Fibonacci server is running at Port: %v \n", shared.CALCULATOR_PORT)

	evolutive.EvolutiveInjector{}.Start("srhtcp", "srhudp", 15*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
