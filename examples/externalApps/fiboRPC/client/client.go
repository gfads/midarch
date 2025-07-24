package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

func clientX(client *rpc.Client) {
	var n, SAMPLE_SIZE, AVERAGE_WAITING_TIME int
	if len(os.Args) >= 2 {
		n, _ = strconv.Atoi(os.Args[1])
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
		AVERAGE_WAITING_TIME = 60
	} else {
		n, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("FIBONACCI_PLACE", "40"))
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValue("SAMPLE_SIZE"))
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(shared.EnvironmentVariableValue("AVERAGE_WAITING_TIME"))
	}
	fmt.Println("dateTime;info;sequential;response_time") //"FILE_SIZE / SAMPLE_SIZE / AVERAGE_WAITING_TIME:", FILE_SIZE, "/", SAMPLE_SIZE, "/", AVERAGE_WAITING_TIME)

	//durations := [SAMPLE_SIZE]time.Duration{}

	rand.Seed(time.Now().UnixNano())
	// invoke remote method
	for x := 0; x < SAMPLE_SIZE; x++ {

		t1 := time.Now()
		//fibo.Fibo(n)
		args := n // Fibonacci place
		var reply int
		err := client.Call("Fibonacci.FiboRPC", args, &reply)
		if err != nil {
			log.Fatal(";error", err, ";;\n")
		}
		//fmt.Printf("Fibo: %d => %d\n", args, reply)
		t2 := time.Now()

		duration := t2.Sub(t1)
		//time.Sleep(3*time.Second)

		//durations[i] = t2.Sub(t1)

		log.Printf(";ok;%d;%f\n", x+1, float64(duration.Nanoseconds())/1000000)

		// Normally distributed waiting time between calls with an average of 60 milliseconds and standard deviation of 20 milliseconds
		var rd = int(math.Round((rand.NormFloat64() + 3) * float64(AVERAGE_WAITING_TIME/3)))
		if rd > 0 {
			time.Sleep(time.Duration(rd) * time.Millisecond)
		}
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
	timeToRun, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("TIME_TO_START_CLIENT", "8"))
	lib.PrintlnDebug("Waiting", timeToRun, "seconds for naming server and server to get up")
	time.Sleep(time.Duration(timeToRun) * time.Second)

	// start configuration
	//frontend.FrontEnd{}.Deploy("midfibonacciclient.madl")

	// proxy to naming service
	//namingProxy := factories.LocateNaming()

	// obtain proxy of fibonacci
	//s := "Fibonacci"
	//proxy1, ok := namingProxy.Lookup(s)
	//if !ok {
	//	fmt.Printf("Client:: Service '%v' not registered in Naming Service!! \n", s)
	//	os.Exit(0)
	//}
	//fibo1 := proxy1.(components.Fibonacciproxy)

	client, err := rpc.Dial("tcp", "server:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("RPC error while dialing:", err)
	}
	clientX(client)
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}
