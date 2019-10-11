package main

import (
	"net/rpc"
	"log"
	"time"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"fmt"
	"strconv"
	"os"
)

func main() {
	var reply int

	// read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// connect to server
	client, err := rpc.DialHTTP("tcp", parameters.NAMING_HOST+":"+strconv.Itoa(parameters.NAMING_PORT))
	if err != nil {
		log.Fatal("Server not ready:", err)
	}

	// make requests
	var t1, t2 time.Time
	var x float64
	for i := 0; i < parameters.SAMPLE_SIZE; i++ {

		// prepare request
		args := shared.Args{A: 38}

		// make request
		t1 = time.Now()
		client.Call("Fibonacci.Fibo", args, &reply)
		t2 = time.Now()

		// print time
		x = float64(t2.Sub(t1).Nanoseconds()) / 1000000
		fmt.Println(x)

		time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	}
	client.Close()
}
