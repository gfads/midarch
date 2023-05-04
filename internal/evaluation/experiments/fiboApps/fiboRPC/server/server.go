package main

import (
	"github.com/gfads/midarch/pkg/apps/fibomiddleware/impl"
	"github.com/gfads/midarch/pkg/shared"
	"log"
	"net"
	"net/rpc"
)

func main() {
	fibonacci := new(impl.Fibonacci)

	rpc.Register(fibonacci)

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	rpc.Accept(ln)
}
