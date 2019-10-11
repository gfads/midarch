package main

import (
	"newsolution/shared/shared"
	"gmidarch/development/frontend"
	"fmt"
)

func main(){

	err := frontend.FrontEnd{}.Deploy("SenderReceiver.madls")
	shared.CheckError(err,"MAIN")

	fmt.Scanln()
}
