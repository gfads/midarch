package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"shared"
	"strconv"
)

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      int
}

var lnSRH net.Listener
var connSRH net.Conn
var firstListenerSRH bool

func NewSRH() SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = "localhost" // TODO
	r.Port = 1313        // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	firstListenerSRH = true

	return *r
}

func (SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := "localhost"            // TODO
	port := shared.CALCULATOR_PORT // TODO

	// create listener
	err := *new(error)
	if firstListenerSRH {
		firstListenerSRH = false
		lnSRH, err = net.Listen("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("SRH:: %v\n",err)
		}
	}

	// accept connections
	connSRH, err = lnSRH.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = connSRH.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msgTemp := make([]byte, sizeInt)
	_, err = connSRH.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp} // TODO
}

func (SRH) I_Send(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, 4)
	l := uint32(len(msgTemp))
	binary.LittleEndian.PutUint32(size, l)
	connSRH.Write(size)

	err := *new(error)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = connSRH.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}
