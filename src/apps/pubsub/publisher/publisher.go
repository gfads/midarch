package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared/factories"
	"strconv"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("queueclient.madl")

	queueingroxy := factories.FactoryQueueing()
	idx := 0

	for {
		msg1 := "msg [" + strconv.Itoa(idx) + "]"
		r := queueingroxy.Publish("Topic01", msg1)
		if !r {
			fmt.Printf("Publisher:: Message not enqueued\n")
		} else {
			fmt.Printf("Publisher:: Message enqueued:: %v\n",msg1)
		}
		idx++
		//time.Sleep(1 * time.Second)
	}
}
