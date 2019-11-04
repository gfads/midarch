package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"injector/evolutive"
)

func main() {
	frontend.FrontEnd{}.Deploy("middlewareserver.madl")

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	inj.Start("marshaller")

	fmt.Scanln()
}