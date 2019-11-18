package main

import (
	"fmt"
	"gmidarch/development/components"
	"gmidarch/execution/frontend"
	"os"
	"shared"
	"shared/factories"
	"time"
)

func main() {
	// start configuration
	frontend.FrontEnd{}.Deploy("midfibonacciclient.madl")

	// proxy to naming service
	namingProxy := factories.LocateNaming()

	// obtain proxy
	s := "Fibonacci"
	f,ok := namingProxy.Lookup(s)
	if !ok {
		fmt.Printf("Client:: Service '%v' not registered in Naming Service!! \n",s)
		os.Exit(0)
	}

	fibo := f.(components.Fibonacciproxy)

	fmt.Printf("Client:: Got a Proxy to Fibonacci\n")

	// invoke remote method
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

		t1 := time.Now()
		fibo.Fibo(38)
		t2 := time.Now()

		x := float64(t2.Sub(t1).Nanoseconds()) / 1000000
		fmt.Printf("%F \n", x)
		//time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	}
}