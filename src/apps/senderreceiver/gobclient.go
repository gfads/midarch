package main

import (
	"encoding/gob"
	"fmt"
	"gmidarch/development/miop"
	"log"
	"net"
	"os"
	"shared"
)

func main() {
	servAddr := "localhost" + ":" + shared.FIBONACCI_PORT // TODO
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	enc := gob.NewEncoder(conn)

	pckt := miop.Packet{}

	pckt.Hdr.Magic = "MIOP"
	err = enc.Encode(&pckt)
	if err != nil {
		fmt.Printf("Error on sending\n")
		os.Exit(0)
	}
}
