package main

import (
	"fmt"
	"injector/evolutive"
	"shared"
)

func main() {

	//elem := "notificationconsumer"
	elem := "fibonacciinvokerm"
	//elem := "receiver"
	inj := evolutive.EvolutiveInjector{}
	go inj.Start(elem, shared.INJECTION_TIME)

	fmt.Scanln()
}
