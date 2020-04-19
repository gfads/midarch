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

type SRHQuic struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

// TODO: Can't be part of SRHQuic? I've to change the names because the scope is the entire package
var ConnsSRHQuic []net.Conn
var LnSRHQuic net.Listener

var c1Quic = make(chan []byte)
var c2Quic = make(chan []byte)
var currentConnectionQuic = -1
var stateQuic = 0

func NewSRHQuic() SRHQuic {

	r := new(SRHQuic)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHQuic) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHQuic).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHQuic).I_Send(msg, info, elemInfo)
	}
}

func (e SRHQuic) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	if LnSRHQuic == nil { // listener was not created yet
		//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		//if err != nil {
		//	log.Fatalf("SRH:: %v\n", err)
		//}
		ln, err := tls.Listen("tcp4", host+":"+port, getServerTLSQuicConfig())
		if err != nil {
			log.Fatalf("SRHQuic:: %v\n", err)
		}
		LnSRHQuic = ln
	}

	switch stateQuic {
	case 0:
		go acceptAndReadQuic(currentConnectionQuic, c1Quic)
		stateQuic = 1
	case 1:
		go acceptAndReadQuic(currentConnectionQuic, c1Quic)
		stateQuic = 2
	case 2:
		go readQuic(currentConnectionQuic, c2Quic)
	}

	//go acceptAndRead(currentConnectionQuic, c1Quic, done)
	//go read(currentConnectionQuic, c2Quic, done)

	select {
	case msgTemp := <-c1Quic:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2Quic:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnectionQuic = nextConnectionQuic()
}

// TODO : can't be part of SRHQuic?
func acceptAndReadQuic(currentConnectionQuic int, c chan []byte) {

	// accept connections
	temp, err := LnSRHQuic.Accept()
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRHQuic = append(ConnsSRHQuic, temp)
	currentConnectionQuic++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRHQuic[currentConnectionQuic]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRHQuic:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

// TODO : can't be part of SRHQuic?
func readQuic(currentConnectionQuic int, c chan []byte) {

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRHQuic[currentConnectionQuic]

	_, err := tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRHQuic:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (e SRHQuic) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnsSRHQuic[currentConnectionQuic].Write(size)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = ConnsSRHQuic[currentConnectionQuic].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRHQuic:: %v\n", err)
		os.Exit(1)
	}
}

// TODO : can't be part of SRHQuic?
func nextConnectionQuic() int {
	r := -1

	if currentConnectionQuic == -1 {
		r = 0
	} else {
		if (currentConnectionQuic + 1) == len(ConnsSRHQuic) {
			r = 0
		} else {
			r = currentConnectionQuic + 1
		}
	}
	return r
}

func getServerTLSQuicConfig() *tls.Config {
	pwd, _ := os.Getwd()
	// TODO: adjust path to crt and key files
	cert, err := tls.LoadX509KeyPair(pwd+"/apps/server/ssl/localhost.crt", pwd+"/apps/server/ssl/localhost.key")
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"MidArchQuic"}, // TODO: Verify what NextProtos should be
	}
	return tlsConfig
}