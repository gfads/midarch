package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("midnamingserver.madl")

	fmt.Scanln()
}