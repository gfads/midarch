package main

import (
	"apps/fibomiddleware/impl"
	"log"
	"net"
	"net/rpc"
	"shared"
)

func main() {
	fibonacci := new(impl.Fibonacci)

	rpc.Register(fibonacci)

	addr, err := net.ResolveTCPAddr("tcp", "localhost:" + shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	rpc.Accept(ln)
}