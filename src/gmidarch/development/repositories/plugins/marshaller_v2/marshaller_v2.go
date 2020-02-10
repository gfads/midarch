package main

import (
	"encoding/json"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
	"shared"
)

type Marshaller struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Gettype() interface{} {
	return Marshaller{}
}

func NewReceiver() Marshaller {

	// create a new instance of client
	r := new(Marshaller)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
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
		fmt.Println("Marshaller:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}
}
