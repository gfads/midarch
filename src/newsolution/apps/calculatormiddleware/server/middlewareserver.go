package main

import (
	"fmt"
	"newsolution/gmidarch/execution/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("middlewareserver.madl")

	fmt.Scanln()
}