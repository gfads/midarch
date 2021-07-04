package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy("senderreceiver.madl","","","","","")

	fmt.Scanln()
}