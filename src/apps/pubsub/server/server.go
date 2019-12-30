package main

import (
	"apps/fibomiddleware/impl"
	"fmt"
	"gmidarch/execution/frontend"
	"os"
	"shared/factories"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("queueclient.madl")

	queueingroxy := factories.FactoryQueueing()
	reqQueue := "requests"
	repQueue := "replies"
	chn := make(chan interface{})

	ok := queueingroxy.Subscribe(reqQueue, chn)


	if ok {
		fmt.Printf("Server:: Server subscribed to queue '%v'\n", reqQueue)
	} else {
		fmt.Printf("Server:: Server not subscribed to queue '%v'\n", reqQueue)
		os.Exit(1)
	}

	for {

		reqMsg := <-chn

		x1 := reqMsg.(map[string]interface{})
		x2 := x1["Msg"].(map[string]interface{})
		//fmt.Printf("Server:: %v\n",x2)
		x3 := x2["Args"].([]interface{})
		p1 := int(x3[0].(float64))
		queueingroxy.Publish(repQueue, impl.Fibonacci{}.F(p1))
		/*if ok {
			fmt.Printf("Server:: Message sent to Client:: %v\n", repMsg)
		} else {
			fmt.Printf("Server:: Message NOT sent to Client:: %v\n", repMsg)
			os.Exit(1)
		}
		*/
	}
}
