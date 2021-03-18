package main

import (
	"fmt"
	"gmidarch/development/components"
	"gmidarch/execution/frontend"
	"os"
	"shared"
	"shared/factories"
	"strconv"
	"time"
)

func clientX(fibo components.HttpProxy){
	var n, SAMPLE_SIZE int
	if len(os.Args) >= 2 {
		n, _ = strconv.Atoi(os.Args[1])
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
	}else{
		n, _ = strconv.Atoi(shared.EnvironmentVariableValue("FIBONACCI_PLACE"))
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValue("SAMPLE_SIZE"))
	}
	fmt.Println("FIBONACCI_PLACE / SAMPLE_SIZE :", n, "/", SAMPLE_SIZE)

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
	// Wait for server to get up
	time.Sleep(10 * time.Second)

	// start configuration
	frontend.FrontEnd{}.Deploy("httpclient.madl")

	// proxy to naming service
	fibo1 := factories.GetHttpProxy("server", shared.HTTP_PORT)

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
