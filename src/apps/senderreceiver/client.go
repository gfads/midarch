package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"shared"
	"strconv"
)

func main() {

	// connect to server
	var err error
	var conn net.Conn

	host := "localhost"
	port := shared.FIBONACCI_PORT

	servAddr := host + ":" + strconv.Itoa(port) // TODO
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	str := "test"
	msgToServer := []byte(str)

	// send size & message
	size := make([]byte, 4)
	l := uint32(len(msgToServer))
	binary.LittleEndian.PutUint32(size, l)

	_, err = conn.Write(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	_, err = conn.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive size & reply
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	msgFromServer := make([]byte,binary.LittleEndian.Uint32(size))
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	fmt.Printf("Message received:: %v\n",msgFromServer)
}
