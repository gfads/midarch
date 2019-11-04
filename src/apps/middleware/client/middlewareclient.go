package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"injector/evolutive"
)

func main() {
	frontend.FrontEnd{}.Deploy("middlewareclient.madl")

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	inj.Start("marshaller")

	fmt.Scanln()
}