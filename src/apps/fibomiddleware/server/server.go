package main

import (
	"fmt"
	"gmidarch/development/components"
	"gmidarch/execution/frontend"
	"injector/evolutive"
	"os"
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
	ok := namingProxy.Register("Fibonacci", fiboProxy)
	if !ok {
		fmt.Printf("Server:: Service 'Fibonacci' already registered in the Naming Server\n")
		os.Exit(0)
	}

	fmt.Printf("Server:: Fibonacci server is running at Port: %v \n",shared.FIBONACCI_PORT)

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	go inj.Start("fibonacciinvokerm")
	fmt.Printf("Server:: Evolutive injector running \n")


	fmt.Scanln()
	fmt.Println("done")
}
