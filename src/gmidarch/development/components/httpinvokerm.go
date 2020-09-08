package components

import (
	"apps/httpserver/impl"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
)

type HttpInvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewHttpInvokerM() HttpInvokerM {
	r := new(HttpInvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e HttpInvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
}

func (HttpInvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) { // TODO
	// unmarshall
	payload := msg.Payload.([]byte)
	request := messages.HttpRequest{}
	request.Unmarshal(payload)

	switch request.Method {
	case "GET":
		response := messages.HttpResponse{}
		impl.RequestListener(request, &response)

		msgTemp := response.Marshal()
		*msg = messages.SAMessage{Payload: msgTemp}
	default:
		fmt.Printf("HttpInvokerM:: Method '%v' not implemented by Http Service\n", request.Method)
		os.Exit(0)
	}
	//fmt.Println("HttpInvokerM.I_Process finished")
}
