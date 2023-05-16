package main

import (
	"fmt"
	"github.com/gfads/midarch/examples/fibonaccidistributed/fibonacciProxy"
	"github.com/gfads/midarch/examples/fibonaccidistributed/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
	"os"
	"sync"
)

func main() {
	// Example setting environment variable MIDARCH_BUSINESS_COMPONENTS_PATH on code, may be set on system environment variables too
	os.Setenv("MIDARCH_BUSINESS_COMPONENTS_PATH",
		shared.DIR_BASE+"/examples/fibonaccidistributed/middleware,"+shared.DIR_BASE+"/examples/fibonaccidistributed/fibonacciProxy")

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["srh"] = messages.EndPoint{Host: "0.0.0.0", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("naming.madl", args, map[string]interface{}{
		"FibonacciInvoker": &middleware.FibonacciInvoker{},
		"FibonacciProxy":   &fibonacciProxy.FibonacciProxy{},
	})

	fmt.Printf("Naming server is running at Port: %v \n", shared.NAMING_PORT)
	//evolutive.EvolutiveInjector{}.Start("srhtcp", 40*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
