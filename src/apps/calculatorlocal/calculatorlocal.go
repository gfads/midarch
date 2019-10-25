package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main() {

	fe := frontend.FrontEnd{}
	fe.Deploy("calculatorlocal.madl")

	fmt.Scanln()
}
