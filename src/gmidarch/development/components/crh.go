package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"strconv"
)

type CRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

//var connCRH net.Conn

func NewCRH() CRH {

	// create a new instance of Server
	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string] net.Conn)

	return *r
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	argsTemp := msg.Payload.([]interface{})
	host := argsTemp[0].(string)
	port := argsTemp[1].(int)
	msgToServer := argsTemp[2].([]byte)

	// connect to server
	var err error
	var key = host+strconv.Itoa(port)
	if _,ok := c.Conns[key]; !ok {
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

	// send message's size
	sizeMsgToServer := make([]byte, 4)
	l := uint32(len(msgToServer))
	binary.LittleEndian.PutUint32(sizeMsgToServer, l)
	c.Conns[key].Write(sizeMsgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = c.Conns[key].Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive message's size
	sizeMsgFromServer := make([]byte, 4)
	_, err = c.Conns[key].Read(sizeMsgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}
	sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

	// receive reply
	msgFromServer := make([]byte, sizeFromServerInt)
	_, err = c.Conns[key].Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}

/*
func (CRH) I_ProcessOld(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	argsTemp := msg.Payload.([]interface{})
	host := argsTemp[0].(string)
	port := argsTemp[1].(int)
	msgToServer := argsTemp[2].([]byte)

	// connect to server
	var err error
	servAddr := host + ":" + strconv.Itoa(port) // TODO
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	connCRH, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Client:: %v\n", err)
	}

	// send message's size
	sizeMsgToServer := make([]byte, 4)
	l := uint32(len(msgToServer))
	binary.LittleEndian.PutUint32(sizeMsgToServer, l)
	connCRH.Write(sizeMsgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = connCRH.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive message's size
	sizeMsgFromServer := make([]byte, 4)
	_, err = connCRH.Read(sizeMsgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}
	sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

	// receive reply
	msgFromServer := make([]byte, sizeFromServerInt)
	_, err = connCRH.Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
*/