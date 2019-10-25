package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"injector/evolutive"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("senderreceiver.madl")

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	inj.Start("receiver")

	fmt.Scanln()
}
