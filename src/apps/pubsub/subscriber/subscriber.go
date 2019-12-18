package main

import (
	"fmt"
	"gmidarch/execution/frontend"
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

	queueing.Subscribe(topic01, chn)

	for {
		fmt.Printf("Subscriber:: %v\n", <- chn)
	}
}
