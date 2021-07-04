package main

import (
	"fmt"
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	"shared"
)

func main() {
	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make (map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host:"localhost",Port:shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("calculatordistributedclientmid.madl",args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host:shared.NAMING_HOST,Port:shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux,ok := namingProxy.Lookup("Calculator")

	calc := calculatorproxy.NewCalculatorProxy()
	if !ok {
		shared.ErrorHandler(shared.GetFunction(),"Service 'Calculator' not found in Naming Service")
	}

	calc = aux.(calculatorproxy.Calculatorproxy)
	fmt.Println(calc.Add(1,2))

	fmt.Scanln()
}