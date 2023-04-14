package main

import (
	"fmt"
	"github.com/gfads/midarch/src/gmidarch/execution/frontend"
	"github.com/gfads/midarch/src/shared"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy("namingclientmid.madl", "localhost", shared.NAMING_PORT) // serverhost, serverport

	fmt.Scanln()
}
