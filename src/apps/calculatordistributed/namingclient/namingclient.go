package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"shared"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy("namingclientmid.madl","localhost",shared.NAMING_PORT) // serverhost, serverport

	fmt.Scanln()
}
