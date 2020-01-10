package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"os"
	"shared"
	"shared/factories"
	"time"
)

func main() {
	durations := [shared.SAMPLE_SIZE]time.Duration{}

	fe := frontend.FrontEnd{}
	fe.Deploy("queueclient.madl")

	queueingroxy := factories.FactoryQueueing()
	reqQueue := "requests"
	repQueue := "replies"
	chn := make(chan interface{})

	ok := queueingroxy.Subscribe(repQueue, chn)
	if ok {
		fmt.Printf("Client:: Client subscribed to queue '%v'\n", repQueue)
	} else {
		fmt.Printf("Client:: Client not subscribed to queue '%v'\n", repQueue)
		os.Exit(1)
	}

	N := 10
	reqMsg := shared.Request{Op: "Fibo", Args: []interface{}{N}}

	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		t1 := time.Now()
		queueingroxy.Publish(reqQueue, reqMsg)
		<-chn
		t2 := time.Now()

		durations[i] = t2.Sub(t1)

		fmt.Printf("[%v] %v\n",i,float64(durations[i].Nanoseconds())/1000000)
		//time.Sleep(10 * time.Millisecond)
	}

	//totalTime := time.Duration(0)
	//for i := range durations {
	//	totalTime += durations[i]
	//}

	//fmt.Printf("Tempo Total [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime)
	//fmt.Printf("Tempo Médio [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime/shared.SAMPLE_SIZE)
}
