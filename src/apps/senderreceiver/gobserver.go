package main

import (
	"encoding/gob"
	"fmt"
	"gmidarch/development/miop"
	"log"
	"net"
	"shared"
)

func main() {

	// create listener if necessary
	servAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatalf("SRH:: %v\n", err)
	}
	ln, err := net.ListenTCP("tcp", servAddr)
	if err != nil {
		log.Fatalf("SRH:: %v\n", err)
	}

	// accept connections
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	dec := gob.NewDecoder(conn)

	pkct := &miop.Packet{}
	err = dec.Decode(pkct)

	fmt.Printf("%v\n", pkct)

	// receive size & message

	// send reply
}
