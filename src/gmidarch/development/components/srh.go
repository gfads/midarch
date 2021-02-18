package components

import (
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

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnsSRH []net.Conn
var LnSRH net.Listener

var c1 = make(chan []byte)
var c2 = make(chan []byte)
var currentConnection = -1
var state = 0

func NewSRH() SRH {

	r := new(SRH)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRH).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRH).I_Send(msg, info, elemInfo)
	}
}

func (e SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "127.0.0.1" // TODO

	if LnSRH == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		LnSRH, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
	}

	switch state {
	case 0:
		go acceptAndRead(currentConnection, c1)
		state = 1
	case 1:
		go read(currentConnection, c1)
		state = 2
	case 2:
		go read(currentConnection, c1)
	}

	//go acceptAndRead(currentConnection, c1, done)
	//go read(currentConnection, c2, done)

	select {
	case msgTemp := <-c1:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2:
		*msg = messages.SAMessage{Payload: msgTemp}
	}

	currentConnection = nextConnection()
}

func acceptAndRead(currentConnection int, c chan []byte) {

	// accept connections
	temp, err := LnSRH.Accept()
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	ConnsSRH = append(ConnsSRH, temp)
	currentConnection++

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRH[currentConnection]
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRH:: Accept and Read\n")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
	c <- msgTemp
}

func read(currentConnection int, c chan []byte) {

	// receive size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := ConnsSRH[currentConnection]

	_, err := tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRH:: Read\n")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	// receive message
	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	c <- msgTemp

	return
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnsSRH[currentConnection].Write(size)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = ConnsSRH[currentConnection].Write(msgTemp)
	if err != nil {
		fmt.Printf("SRH:: %v\n", err)
		os.Exit(1)
	}
}

func nextConnection() int {
	r := -1

	if currentConnection == -1 {
		r = 0
	} else {
		if (currentConnection + 1) == len(ConnsSRH) {
			r = 0
		} else {
			r = currentConnection + 1
		}
	}
	return r
}
