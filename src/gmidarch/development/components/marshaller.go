package components

import (
	"fmt"
	graphs2 "gmidarch/development/artefacts/graphs"
	impl2 "gmidarch/development/impl"
	messages2 "gmidarch/development/messages"
	miop2 "gmidarch/development/miop"
	"os"
	shared2 "shared"
)

type Marshaller struct {
	CSP       string
	Graph     graphs2.ExecGraph
	Behaviour string
}

func NewMarshaller() Marshaller {

	// create a new instance of Server
	r := new(Marshaller)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (Marshaller) I_Process(msg *messages2.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared2.Request)
	op := req.Op

	switch op {
	case "marshall":
		p1 := req.Args[0].(miop2.Packet)
		r := impl2.MarshallerImpl{}.Marshall(p1)
		*msg = messages2.SAMessage{Payload: r}
	case "unmarshall":
		p1 := req.Args[0].([]byte)
		r := impl2.MarshallerImpl{}.Unmarshall(p1)
		*msg = messages2.SAMessage{Payload: r}
	default:
		fmt.Println("Marshaller:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}
}
