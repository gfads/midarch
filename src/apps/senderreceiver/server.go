package main

import (
	"encoding/binary"
	"log"
	"net"
	"shared"
	"strconv"
)

func main() {
	host := "localhost"
	port := shared.FIBONACCI_PORT
	var ln net.Listener
	var conn net.Conn

	// create listener
	servAddr, err := net.ResolveTCPAddr("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("SRH:: %v\n", err)
	}

	ln, err = net.ListenTCP("tcp", servAddr)
	if err != nil {
		log.Fatalf("SRH:: %v\n", err)
	}

	// accept connections
	conn, err = ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive size & message
	size := make([]byte, 4)
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	msgFromClient := make([]byte,binary.LittleEndian.Uint32(size))

	_, err = conn.Read(msgFromClient)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send size & message
	msgToClient := msgFromClient

	l := uint32(len(msgToClient))
	binary.LittleEndian.PutUint32(size, l)
	_, err = conn.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	_, err = conn.Write(msgToClient)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}
