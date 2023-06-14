package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gfads/midarch/examples/calculatordistributed/externalcomponents"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	evolutive "github.com/gfads/midarch/pkg/injector"
	"github.com/gfads/midarch/pkg/shared"
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
	fe.Deploy(frontend.DeployOptions{FileName: "calculatordistributedservermid.madl", Args: args})

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	// Create proxy to calculatorimpl
	calcProxy := externalcomponents.NewCalculatorProxy(generic.ProxyConfig{
		Host: shared.CALCULATOR_HOST,
		Port: shared.CALCULATOR_PORT,
	})

	// Register calculatorimpl in Lookup
	ok := namingProxy.Register("Calculator", calcProxy)

	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "'Calculator' already registered in the Naming Server")
	}

	fmt.Printf("Calculator server is running at Port: %v \n", shared.CALCULATOR_PORT)

	evolutive.EvolutiveInjector{}.Start("srhtcp", "srhudp", 15*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
