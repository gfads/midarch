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

		// receive request
		reqMsg := <-chn
		//fmt.Printf("Server:: Receive request:: [%v] \n",reqMsg)
		x1 := reqMsg.(map[string]interface{})
		x2 := x1["Msg"].(map[string]interface{})
		x3 := x2["Args"].([]interface{})
		p1 := int(x3[0].(float64))

		// calculate and publish result
		queueingroxy.Publish(repQueue, impl.Fibonacci{}.F(p1))
		//fmt.Printf("Server:: Published response\n")
	}
}
