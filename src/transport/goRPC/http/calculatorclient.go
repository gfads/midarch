package main

import (
	"log"
	"fmt"
	"time"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"strconv"
	"net/rpc"
)

func main() {

	var reply int
	// connect to server
	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(parameters.CALCULATOR_PORT))
	if err != nil {
		log.Fatal("Server not ready:", err)
	}

	// make requests
	for i := 0; i < parameters.SAMPLE_SIZE; i++ {

		args := shared.Args{A: i, B: i}

		t1 := time.Now()
		client.Call("Calculator.Add", args, &reply)
		t2 := time.Now()

		x := float64(t2.Sub(t1).Nanoseconds()) / 1000000
		fmt.Println(x)
	}
}
