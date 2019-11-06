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

type CRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

//var connCRH net.Conn

func NewCRH() CRH {

	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]net.Conn)

	return *r
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	argsTemp := msg.Payload.([]interface{})
	host := argsTemp[0].(string)         // host
	port := argsTemp[1].(int)            // port
	msgToServer := argsTemp[2].([]byte)  // message

	// connect to server
	var err error
	var key = host + strconv.Itoa(port)
	if _, ok := c.Conns[key]; !ok { // no connection open yet
		servAddr := host + ":" + strconv.Itoa(port) // TODO
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			log.Fatalf("Client:: %v\n", err)
		}

		c.Conns[key], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatalf("Client:: %v\n", err)
		}
	}

	// configure conn to be used
	conn := c.Conns[key]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}