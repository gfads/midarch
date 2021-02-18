package main

import (
	"fmt"
	"gmidarch/development/components"
	"gmidarch/execution/frontend"
	"shared"
	"shared/factories"
	"time"
)

func clientX(fibo components.Http2Proxy){

	//n,_ := strconv.Atoi(os.Args[1])
	//SAMPLE_SIZE,_ := strconv.Atoi(os.Args[2])
	n := 10
	SAMPLE_SIZE := 10

	//durations := [SAMPLE_SIZE]time.Duration{}

	// invoke remote method
	for i := 0; i < SAMPLE_SIZE; i++ {

		t1 := time.Now()
		fibo.Fibo(n)
		t2 := time.Now()

		duration := t2.Sub(t1)

		//durations[i] = t2.Sub(t1)

		fmt.Printf("%v\n",float64(duration.Nanoseconds())/1000000)
	}

	//totalTime := time.Duration(0)
	//for i := range durations {
	//	totalTime += durations[i]
	//}

	//fmt.Printf("Tempo Total [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime)
	//fmt.Printf("Tempo Médio [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime/shared.SAMPLE_SIZE)
}

func main() {
	// start configuration
	frontend.FrontEnd{}.Deploy("http2client.madl")

	// proxy to naming service
	fibo1 := factories.GetHttp2Proxy("https://localhost", shared.HTTP_PORT)

	//// obtain proxy of fibonacci
	//s := "Fibonacci"
	//proxy1, ok := namingProxy.Lookup(s)
	//if !ok {
	//	fmt.Printf("Client:: Service '%v' not registered in Naming Service!! \n", s)
	//	os.Exit(0)
	//}
	//fibo1 := proxy1.(components.Fibonacciproxy)

	clientX(fibo1)
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}
