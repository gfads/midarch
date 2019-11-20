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
	proxy,ok := namingProxy.Lookup(s)
	if !ok {
		fmt.Printf("Client:: Service '%v' not registered in Naming Service!! \n",s)
		os.Exit(0)
	}

	fibo := proxy.(components.Fibonacciproxy)

	durations := [5000]time.Duration{}

	// invoke remote method
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

		t1 := time.Now()
		fibo.Fibo(10)
		t2 := time.Now()

		durations[i] = t2.Sub(t1)

		//if i >= 1500 {
		//	fmt.Println(float64(t2.Sub(t1).Nanoseconds())/float64(1000000))
		//}

		//x := float64(t2.Sub(t1).Nanoseconds()) / 1000000
		//fmt.Printf("%v \n", x)
		//time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	}

	totalTime := time.Duration(0)
	for i := range durations{
		totalTime += durations[i]
	}

	fmt.Printf("Tempo Total: %v\n",totalTime)
	fmt.Printf("Tempo Médio: %v\n",totalTime/5000)
}