package main

import (
	"net/rpc"
	"net"
	"log"
	"fmt"
	"net/http"
	"newsolution/shared/parameters"
	"strconv"
	"apps/calculator/implrpc"
	"newsolution/shared/shared"
)

func main() {

	// create new instance of calculator
	calculator := new(implrpc.Calculator)

	// create new goRPC server
	server := rpc.NewServer()
	server.RegisterName("Calculator", calculator)

	// associate a http handler to server
	server.HandleHTTP("/", "/debug")

	// create tcp listen
	l, e := net.Listen("tcp", shared.ResolveHostIp()+":"+strconv.Itoa(parameters.CALCULATOR_PORT))
	if e != nil {
		log.Fatal("Server not started:", e)
	}

	// wait for calls
	fmt.Println("Server is running... at "+shared.ResolveHostIp()+" Port "+strconv.Itoa(parameters.CALCULATOR_PORT))
	http.Serve(l, nil)
}


