package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
)

func main() {
	frontend.FrontEnd{}.Deploy("midnamingserver.madl")

	fmt.Printf("Naming server ready at port '%v' ...\n",shared.NAMING_PORT)

	fmt.Scanln()
}