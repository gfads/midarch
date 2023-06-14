package main

import (
	"fmt"
	"time"

	"github.com/gfads/midarch/examples/calculatordistributed/externalcomponents"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
)

func main() {
	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: "localhost", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy(frontend.DeployOptions{
		FileName: "calculatordistributedclientmid.madl",
		Args:     args,
		Components: map[string]interface{}{
			"Calculatorinvoker": &externalcomponents.Calculatorinvoker{},
			"Calculatorproxy":   &externalcomponents.Calculatorproxy{},
		}})

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux, ok := namingProxy.Lookup("Calculator")
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Service 'Calculator' not found in Naming Service")
	}

	calc := aux.(*externalcomponents.Calculatorproxy)
	for x := 0; x < 1000; x++ {
		fmt.Println("Result:", calc.Add(x, 1))
		time.Sleep(200 * time.Millisecond)
	}

	//fmt.Scanln()
	//var wg sync.WaitGroup
	//wg.Add(1)
	//wg.Wait()
}
