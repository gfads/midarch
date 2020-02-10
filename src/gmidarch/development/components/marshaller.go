package components

import (
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
	"shared"
)

type Marshaller struct {
	CSP       string
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewMarshaller() Marshaller {

	r := new(Marshaller)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (Marshaller) I_ProcessMessagePack(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op

	switch op {
	case "marshall":
		p1 := req.Args[0].(miop.Packet)
		r, err := msgpack.Marshal(p1)
		if err != nil {
			log.Fatalf("Marshaller:: Marshall:: %s", err)
		}
		*msg = messages.SAMessage{Payload: r}
	case "unmarshall":
		p1 := req.Args[0].([]byte)
		r := miop.Packet{}
		err := msgpack.Unmarshal(p1,&r)
		if err != nil {
			log.Fatalf("Marshaller:: Unmarshall:: %s", err)
		}
		*msg = messages.SAMessage{Payload: r}
	default:
		fmt.Printf("Marshaller:: Operation '%v' not supported!",op)
		os.Exit(0)
	}
}

func (Marshaller) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op

	switch op {
	case "marshall":
		p1 := req.Args[0].(miop.Packet)
		r, err := json.Marshal(p1)
		if err != nil {
			log.Fatalf("Marshaller:: Marshall:: %s", err)
		}
		*msg = messages.SAMessage{Payload: r}
	case "unmarshall":
		p1 := req.Args[0].([]byte)
		r := miop.Packet{}
		err := json.Unmarshal(p1, &r)
		if err != nil {
			log.Fatalf("Marshaller:: Unmarshall:: %s", err)
		}
		*msg = messages.SAMessage{Payload: r}
	default:
		fmt.Printf("Marshaller:: Operation '%v' not supported!",op)
		os.Exit(0)
	}
}
