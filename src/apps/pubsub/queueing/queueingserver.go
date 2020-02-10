package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main(){
	frontend.FrontEnd{}.Deploy("queueingserver.madl")

	fmt.Printf("Queue server ready!!\n")

	fmt.Scanln()
}
