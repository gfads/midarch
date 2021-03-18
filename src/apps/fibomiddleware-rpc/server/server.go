package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	wg.Add(1)

	// start configuration
	frontend.FrontEnd{}.Deploy("midfibonacciserver-rpc.madl")

	//// proxy to naming service
	//namingProxy := factories.LocateNaming()
	//
	//// register
	//fiboProxy := components.Fibonacciproxy{Host:"server",Port:shared.FIBONACCI_PORT} //shared.ResolveHostIp()
	//ok := namingProxy.Register("Fibonacci", fiboProxy)
	//if !ok {
	//	fmt.Printf("Server:: Service 'Fibonacci' already registered in the Naming Server\n")
	//	os.Exit(0)
	//}

	fmt.Printf("Server:: Fibonacci server is running at Port: %v \n",shared.FIBONACCI_PORT)

	//fmt.Scanln()
	wg.Wait()
	fmt.Println("done")
}
