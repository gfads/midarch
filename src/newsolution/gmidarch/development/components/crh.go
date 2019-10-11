package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"net"
	"strconv"
	"encoding/binary"
	"log"
	"newsolution/gmidarch/development/element"
	"newsolution/shared/shared"
)

type CRH struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewCRH() CRH {

	// create a new instance of Server
	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (c *CRH) Configure(invP, terP *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *graphs.NewExecGraph(3)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(2, 0, newEdgeInfo)

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
