package main

import (
	"fmt"
	"github.com/gfads/midarch/src/gmidarch/execution/frontend"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy("calculatorlocal.madl")

	fmt.Scanln()
}
