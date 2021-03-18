package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	frontend.FrontEnd{}.Deploy("midnamingserver.madl")

	fmt.Printf("Naming server ready at port '%v' ...\n",shared.NAMING_PORT)

	//fmt.Scanln()
	wg.Wait()
}