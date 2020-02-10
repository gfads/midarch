package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"injector/evolutive"
)

func main() {

	// start frontend
	fe := frontend.FrontEnd{}
	fe.Deploy("calculatorlocal.madl")

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	inj.Start("calculatorclient")

	fmt.Scanln()
}
