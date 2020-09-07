package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
)

func main(){

	// start configuration
	frontend.FrontEnd{}.Deploy("httpserver.madl") // Todo: create a port for HttpServer and configure correctly in madl

	//// proxy to naming service
	//namingProxy := factories.LocateNaming()
	//
	//// register
	//// Todo: Create a HttpProxy
	//fiboProxy := components.Fibonacciproxy{Host:shared.ResolveHostIp(),Port:shared.FIBONACCI_PORT}
	//ok := namingProxy.Register("HttpServer", fiboProxy)
	//if !ok {
	//	fmt.Printf("Server:: Service 'HttpServer' already registered in the Naming Server\n")
	//	os.Exit(0)
	//}

	// Todo: create a port for HttpServer and use it correctly in logs
	fmt.Printf("Server:: Http Server is running at Port: %v \n",shared.FIBONACCI_PORT)

	fmt.Scanln()
	fmt.Println("done")
}
