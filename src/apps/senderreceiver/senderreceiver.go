package main

import (
	"fmt"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
)

func main() {
	fe := frontend.NewFrontend()

	args := make (map[string]messages.EndPoint)

	fe.Deploy("senderreceiver.madl", args)

	fmt.Scanln()
}