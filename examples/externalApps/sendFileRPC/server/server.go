package main

import (
	"log"
	"net"
	"net/rpc"

	sendFileImpl "github.com/gfads/midarch/examples/sendfiledistributed/sendfileImpl"
	"github.com/gfads/midarch/pkg/shared"
)

func main() {
	sendFile := new(sendFileImpl.SendFile)

	rpc.Register(sendFile)

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	rpc.Accept(ln)
}
