package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("middlewareserver.madl")

	// Start evolutive injector
	//inj := evolutive.EvolutiveInjector{}
	//inj.Start("receiver")

	fmt.Scanln()
}