package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("clientserver.madl")

	fmt.Scanln()
}