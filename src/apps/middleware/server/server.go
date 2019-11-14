package main

import (
	"fmt"
	"gmidarch/development/components"
	"gmidarch/execution/frontend"
	"shared"
	"shared/factories"
)

func main(){

	// start configuration
	frontend.FrontEnd{}.Deploy("midfibonacciserver.madl")

	// proxy to naming service
	namingProxy := factories.LocateNaming()

	// register
	fiboProxy := components.Fibonacciproxy{Host:shared.ResolveHostIp(),Port:shared.FIBONACCI_PORT}
	namingProxy.Register("Fibonacci", fiboProxy)

//	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fiboProxy.Port))

	fmt.Scanln()
	fmt.Println("done")
}
