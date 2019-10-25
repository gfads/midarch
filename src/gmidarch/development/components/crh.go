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
}

func NewCRH() CRH {

	// create a new instance of Server
	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	argsTemp := msg.Payload.([]interface{})
	host := argsTemp[0].(string)
	port := argsTemp[1].(int)
	msgToServer := argsTemp[2].([]byte)

	// connect to server
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
		if err == nil {
			break
		}

	}

	defer conn.Close()

	// send message's size
	sizeMsgToServer := make([]byte, 4)
	l := uint32(len(msgToServer))
	binary.LittleEndian.PutUint32(sizeMsgToServer, l)
	conn.Write(sizeMsgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive message's size
	sizeMsgFromServer := make([]byte, 4)
	_, err = conn.Read(sizeMsgFromServer)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

	// receive reply
	msgFromServer := make([]byte, sizeFromServerInt)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
