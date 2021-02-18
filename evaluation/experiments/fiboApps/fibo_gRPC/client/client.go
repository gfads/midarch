package main

import (
	fibonacci "apps/fiboApps/fibo_gRPC/proto"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"shared"
	"strconv"
	"time"
)

func clientX(client fibonacci.FibonacciServiceClient){

	n,_ := strconv.Atoi(os.Args[1])
	SAMPLE_SIZE,_ := strconv.Atoi(os.Args[2])
	//n := 38
	//SAMPLE_SIZE := 10

	//durations := [SAMPLE_SIZE]time.Duration{}

	// invoke remote method
	for i := 0; i < SAMPLE_SIZE; i++ {

		t1 := time.Now()
		_, err := client.Fibo(context.Background(), &fibonacci.Request{Place: int32(n)})
		if err != nil {
			log.Fatal("Fibo error:", err)
		}
		//fmt.Printf("Fibo: %d => %d\n", n, response)
		t2 := time.Now()

		duration := t2.Sub(t1)
		//time.Sleep(3*time.Second)

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

	//var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:"+shared.FIBONACCI_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal("gRPC error while dialing:", err)
	}
	defer conn.Close()
	client := fibonacci.NewFibonacciServiceClient(conn)

	clientX(client)
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}