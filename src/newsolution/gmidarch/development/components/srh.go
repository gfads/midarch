package components

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"strconv"
)

type SRH struct {
	Behaviour   string
	Graph graphs.ExecGraph
	Host  string
	Port  int
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

func (s *SRH) Configure(invR, terR *chan messages.SAMessage) {

	// configure the state machine
	s.Graph = *graphs.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)
	msg := new(messages.SAMessage)

	info1 := make([]*interface{}, 3) // host, port, msg
	info1[0] = new(interface{})
	info1[1] = new(interface{})

	*info1[0] = msg
	*info1[1] = make([] interface{}, 3)
	args1HostPort := make([]interface{}, 2)
	args1HostPort[0] = "localhost"
	args1HostPort[1] = 1313
	*info1[1] = args1HostPort

	newEdgeInfo := graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Receive", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info1}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	s.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	s.Graph.AddEdge(2, 3, newEdgeInfo)

	info2 := make([]*interface{}, 1)
	info2[0] = new(interface{})
	*info2[0] = msg
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Send", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info2}
	s.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := "localhost"                // TODO
	port := parameters.CALCULATOR_PORT // TODO

	// create listener
	ln, err = net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("SRH:: %s", err)
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
	conn.Close()
	ln.Close()
}

func (SRH) I_Test1(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Teste 1"}
	*info[0] = 3
	fmt.Printf("SRH:: %v\n", *msg)
}

func (SRH) I_Test2(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Teste 2"}
	*info[0] = 13
	fmt.Printf("SRH:: %v\n", *msg)
}
