package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"shared"
)

type CRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

func NewCRH() CRH {

	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte) // message

	// connect to server
	var err error

	key := host + port

	if _, ok := c.Conns[key]; !ok { // no connection open yet
		servAddr := host + ":" + port // TODO
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			log.Fatalf("Client:: %v\n", err)
		}

		c.Conns[key], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatalf("Client:: %v\n", err)
		}
	}

	conn := c.Conns[key]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
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
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
