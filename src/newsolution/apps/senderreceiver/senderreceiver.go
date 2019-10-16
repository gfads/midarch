package main

import (
	"fmt"
	"newsolution/gmidarch/execution/frontend"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("senderreceiver.madl")

	fmt.Scanln()
}