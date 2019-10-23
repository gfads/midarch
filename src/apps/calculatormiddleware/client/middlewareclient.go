package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("middlewareclient.madl")

	fmt.Scanln()
}