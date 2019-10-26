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

var ln net.Listener
var conn net.Conn
var err error

func NewSRH() SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = "localhost" // TODO
	r.Port = 1313        // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := "localhost"             // TODO
	port := shared.CALCULATOR_PORT  // TODO

	// create listener
	for {
		ln, err = net.Listen("tcp", host+":"+strconv.Itoa(port))
		if err == nil {
			break
		}
	}
	// accept connections
	conn, err = ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msgTemp := make([]byte, sizeInt)
	_, err = conn.Read(msgTemp)
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
	conn.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// close connection
	err = conn.Close()
	if err != nil {
		log.Fatalf("SRH:: %v\n",err)
	}
	ln.Close()
	if err != nil {
		log.Fatalf("SRH:: %v\n",err)
	}

}