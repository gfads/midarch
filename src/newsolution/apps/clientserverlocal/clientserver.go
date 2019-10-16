package main

import (
	"fmt"
	"newsolution/gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("clientserver.madl")

	fmt.Scanln()
}