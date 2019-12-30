package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"os"
	"shared/factories"
)

func main() {

	// Start configuration
	fe := frontend.FrontEnd{}
	fe.Deploy("queueclient.madl")

	// Obtaing proxy to queueing service
	queueing := factories.FactoryQueueing()
	topic01 := "Topic01"
	chn := make(chan interface{})

	if !queueing.Subscribe(topic01, chn){
		fmt.Printf("Subsecriber:: Subscription to topic '%v' failed!!\n",topic01)
		os.Exit(0)
	} else {
		fmt.Printf("Subsecriber:: Subscription to topic '%v' succedded!!\n",topic01)
	}

	for {
		fmt.Printf("Subscriber:: %v\n", <- chn)
	}
}
