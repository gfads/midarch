package main

import (
	"fmt"
	"gmidarch/development/components/proxies/fibonacciProxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	"shared"
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
	fe.Deploy("FibonacciDistributedClientMid.madl", args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux, ok := namingProxy.Lookup("Fibonacci")
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Service 'Fibonacci' not found in Naming Service")
	}

	fibonacci := aux.(*fibonacciProxy.FibonacciProxy)
	for x := 0; x < 1000; x++ {
		fmt.Println("Result:", fibonacci.F(11))
		time.Sleep(200 * time.Millisecond)
	}

	//fmt.Scanln()
	//var wg sync.WaitGroup
	//wg.Add(1)
	//wg.Wait()
}
