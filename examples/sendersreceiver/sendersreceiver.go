package main

import (
	"fmt"

	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy(frontend.DeployOptions{FileName: "sendersreceiver.madl"})

	fmt.Scanln()
}