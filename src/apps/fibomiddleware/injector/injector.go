package main

import (
	"fmt"
	"injector/evolutive"
)

func main(){
	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	go inj.Start("fibonacciinvokerm")

	fmt.Scanln()
}
