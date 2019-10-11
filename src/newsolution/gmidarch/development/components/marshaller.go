package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/impl"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/development/miop"
	"newsolution/shared/shared"
	"os"
)

type Marshaller struct {
	CSP   string
	Graph graphs.ExecGraph
}

func NewMarshaller() Marshaller {

	// create a new instance of Server
	r := new(Marshaller)

	return *r
}

func (m *Marshaller) Configure(invP, terP *chan messages.SAMessage) {

	// configure the state machine
	m.Graph = *graphs.NewExecGraph(3)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	m.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	m.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	m.Graph.AddEdge(2, 0, newEdgeInfo)

}

func (Marshaller) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op

	switch op {
	case "marshall":
		p1 := req.Args[0].(miop.Packet)
		r := impl.MarshallerImpl{}.Marshall(p1)
		*msg = messages.SAMessage{Payload: r}
	case "unmarshall":
		p1 := req.Args[0].([]byte)
		r := impl.MarshallerImpl{}.Unmarshall(p1)
		*msg = messages.SAMessage{Payload: r}
	default:
		fmt.Println("Marshaller:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}
}
