package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("senderreceiver.madl")

	fmt.Scanln()
}