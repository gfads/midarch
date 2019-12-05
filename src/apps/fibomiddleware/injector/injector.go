package main

import (
	"fmt"
	"injector/evolutive"
	"time"
)

func main() {

	elem := "fibonacciinvokerm"
	//elem := "receiver"
	inj := evolutive.EvolutiveInjector{}
	go inj.Start(elem, 1*time.Second)

	fmt.Scanln()
}
