package main

import (
	"fmt"
	"github.com/gfads/midarch/examples/fibonaccidistributed/fibonacciProxy"
	"github.com/gfads/midarch/examples/fibonaccidistributed/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	evolutive "github.com/gfads/midarch/pkg/injector"
	"github.com/gfads/midarch/pkg/shared"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Wait for namingserver to get up
	time.Sleep(8 * time.Second)
	
	// Example setting environment variable MIDARCH_BUSINESS_COMPONENTS_PATH on code, may be set on system environment variables too
	os.Setenv("MIDARCH_BUSINESS_COMPONENTS_PATH",
		shared.DIR_BASE+"/examples/fibonaccidistributed/middleware,"+shared.DIR_BASE+"/examples/fibonaccidistributed/fibonacciProxy")

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	args["srh"] = messages.EndPoint{Host: "0.0.0.0", Port: shared.CALCULATOR_PORT}

	// Deploy configuration
	fe.Deploy("FibonacciDistributedServerMid.madl", args, map[string]interface{}{
		"FibonacciInvoker": &middleware.FibonacciInvoker{},
		"FibonacciProxy":   &fibonacciProxy.FibonacciProxy{},
	})

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

	intervalBetweenInjections, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("INJECTION_INTERVAL", "90"))
	evolutive.EvolutiveInjector{}.StartEvolutiveProtocolInjection("srhtcp", "srhudp", time.Duration(intervalBetweenInjections)*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
