package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// start configuration
	frontend.FrontEnd{}.Deploy("http2server.madl")

	//// proxy to naming service
	//namingProxy := factories.LocateNaming()
	//
	//// register
	//// Todo: Create a HttpProxy
	//fiboProxy := components.Fibonacciproxy{Host:shared.ResolveHostIp(),Port:shared.HTTP_PORT}
	//ok := namingProxy.Register("HttpServer", fiboProxy)
	//if !ok {
	//	fmt.Printf("Server:: Service 'HttpServer' already registered in the Naming Server\n")
	//	os.Exit(0)
	//}

	fmt.Printf("Server:: Http2 Server is running at Port: %v \n",shared.HTTP_PORT)

	//fmt.Scanln()
	wg.Wait()
	fmt.Println("done")
}