package main

import (
	"fmt"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
)

func main() {
	fe := frontend.NewFrontend()

	fe.Deploy("namingclientmid.madl", "localhost", shared.NAMING_PORT) // serverhost, serverport

	fmt.Scanln()
}
