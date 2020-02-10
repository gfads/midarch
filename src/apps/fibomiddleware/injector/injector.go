package main

import (
	"fmt"
	"injector/evolutive"
	"os"
	"shared"
	"strconv"
	"time"
)

func main() {
	var timeBetweenInjections time.Duration
	t1,_ := strconv.Atoi(os.Args[1])
	timeBetweenInjections = time.Duration(t1)

	//elem := "notificationconsumer"
	elem := "fibonacciinvokerm"
	//elem := "receiver"
	inj := evolutive.EvolutiveInjector{}
	shared.INJECTION_TIME = timeBetweenInjections * time.Second
	go inj.Start(elem, shared.INJECTION_TIME)

	fmt.Scanln()
}
