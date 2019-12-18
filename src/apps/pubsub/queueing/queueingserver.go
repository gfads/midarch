package main

import (
	"fmt"
	"gmidarch/execution/frontend"
)

func main(){
	frontend.FrontEnd{}.Deploy("queueserver.madl")

	fmt.Printf("Queue server ready!!\n")

	fmt.Scanln()
}
