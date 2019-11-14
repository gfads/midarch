package components

import (
	"encoding/binary"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"net"
	"os"
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

func (CRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	elem.(CRH).I_Process(msg, info)
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte) // message

	// connect to server
	key := host + ":"+port
	if _, ok := c.Conns[key]; !ok { // no connection open yet
		servAddr := key // TODO
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
	_, err := conn.Write(size)
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

// only used in I_ProcessEncDec
var encCRH msgpack.Encoder
var decCRH msgpack.Decoder

// Use of Encode/Decode - inefficient
func (c CRH) I_ProcessEncDec(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := &miop.Packet{}
	*msgToServer = payload[2].(miop.Packet) // message

	// connect to server
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
		encCRH = *msgpack.NewEncoder(c.Conns[key])
		decCRH = *msgpack.NewDecoder(c.Conns[key])
	}

	err := encCRH.Encode(&msgToServer)
	if err != nil {
		fmt.Printf("CRH:: %v\n",err)
		os.Exit(0)
	}

	msgFromServer := &miop.Packet{}

	err = decCRH.Decode(&msgFromServer)
	if err != nil {
		fmt.Printf("CRH:: %v\n",err)
		os.Exit(0)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
