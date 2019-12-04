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

var N int = 1

func main() {
	// start configuration
	frontend.FrontEnd{}.Deploy("midfibonacciclient.madl")

	// proxy to naming service
	namingProxy := factories.LocateNaming()

	// obtain proxy of fibonacci
	s := "Fibonacci"
	proxy, ok := namingProxy.Lookup(s)
	if !ok {
		fmt.Printf("Client:: Service '%v' not registered in Naming Service!! \n", s)
		os.Exit(0)
	}
	fibo := proxy.(components.Fibonacciproxy)

	durations := [shared.SAMPLE_SIZE]time.Duration{}

	// invoke remote method
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

		t1 := time.Now()
		fibo.Fibo(N)
		t2 := time.Now()

		durations[i] = t2.Sub(t1)

		//time.Sleep(10 * time.Millisecond)
		//fmt.Printf("%v\n",float64(durations[i].Nanoseconds())/1000000)
	}

	totalTime := time.Duration(0)
	for i := range durations {
		totalTime += durations[i]
	}

	fmt.Printf("Tempo Total [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE,totalTime)
	fmt.Printf("Tempo Médio [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE,totalTime/shared.SAMPLE_SIZE)

	fmt.Scanln()
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}