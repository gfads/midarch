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

var ConnSRH net.Conn
var LnSRH net.Listener
var c1 = make(chan []byte)
var c2 = make(chan []byte)


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

	//var err error
	if LnSRH == nil { // listener was not created yet
		tempPort := *elemInfo[0]
		port := tempPort.(string)
		host := "localhost"

		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		LnSRH, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
	}

	// it allows to read/accept simultaneously
	go acceptAndRead(&ConnSRH, LnSRH, c1)
	if ConnSRH != nil {
		go read(ConnSRH, c2)
	}

	select {
	case msgTemp := <-c1:
		*msg = messages.SAMessage{Payload: msgTemp}
	case msgTemp := <-c2:
		*msg = messages.SAMessage{Payload: msgTemp}
	}
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnSRH.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = ConnSRH.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}

func acceptAndRead(conn *net.Conn, ln net.Listener, c chan []byte) {

	// accept connections
	var err error
	*conn, err = ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive size & message
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := *conn
	_, err = tempConn.Read(size)
	if err == io.EOF {
		{
			fmt.Printf("SRH:: Accept and Read")
			os.Exit(0)
		}
	} else if err != nil && err != io.EOF {
		log.Fatalf("SRH::: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	c <- msgTemp
}

func read(conn net.Conn, c chan []byte) {

	// receive size & message
	var err error
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	tempConn := conn
	_, err = tempConn.Read(size)
	if err == io.EOF {
		fmt.Printf("SRH:: Read")
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		log.Fatalf("SRH::: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = tempConn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	c <- msgTemp
}
