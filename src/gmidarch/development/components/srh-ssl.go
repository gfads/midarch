package components

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"os"
	"shared"
)

type SRHSsl struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnsSRHSsl []net.Conn
var LnSRHSsl net.Listener

var c1Ssl = make(chan []byte)
var c2Ssl = make(chan []byte)
var currentConnectionSsl = -1
var stateSsl = 0

func NewSRHSsl() SRHSsl {

	r := new(SRHSsl)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHSsl) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHSsl).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHSsl).I_Send(msg, info, elemInfo)
	}
}

func (e SRHSsl) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	if LnSRHSsl == nil { // listener was not created yet
		//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		//if err != nil {
		//	log.Fatalf("SRH:: %v\n", err)
		//}
		ln, err := tls.Listen("tcp4", host+":"+port, getServerTLSConfig())
		if err != nil {
			log.Fatalf("SRHSsl:: %v\n", err)
		}
		LnSRHSsl = ln
	}

	switch stateSsl {
	case 0:
		go acceptAndReadSsl(currentConnectionSsl, c1Ssl)
		stateSsl = 1
	case 1:
		go readSsl(currentConnectionSsl, c2Ssl)
		stateSsl = 2
	case 2:
		go readSsl(currentConnectionSsl, c2Ssl)
	}

	//go acceptAndRead(currentConnectionSsl, c1Ssl, done)
	//go read(currentConnectionSsl, c2Ssl, done)

	select {
	case msgTemp := <-c1Ssl:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Ssl:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnectionSsl = nextConnectionSsl()
}

func acceptAndReadSsl(currentConnectionSsl int, c chan []byte) {

	// accept connections
	temp, err := LnSRHSsl.Accept()
	if err != nil {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRHSsl = append(ConnsSRHSsl, temp)
	currentConnectionSsl++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRHSsl[currentConnectionSsl]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRHSsl:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

func readSsl(currentConnectionSsl int, c chan []byte) {

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRHSsl[currentConnectionSsl]

	_, err := tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRHSsl:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (e SRHSsl) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnsSRHSsl[currentConnectionSsl].Write(size)
	if err != nil {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = ConnsSRHSsl[currentConnectionSsl].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRHSsl:: %v\n", err)
		os.Exit(1)
	}
}

func nextConnectionSsl() int {
	r := -1

	if currentConnectionSsl == -1 {
		r = 0
	} else {
		if (currentConnectionSsl + 1) == len(ConnsSRHSsl) {
			r = 0
		} else {
			r = currentConnectionSsl + 1
		}
	}
	return r
}

func getServerTLSConfig() *tls.Config {
	if shared.CRT_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'CRT_PATH' not configured\n")
	}

	if shared.KEY_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'KEY_PATH' not configured\n")
	}

	cert, err := tls.LoadX509KeyPair(shared.CRT_PATH, shared.KEY_PATH)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"exemplo"},
	}
	return tlsConfig
}