package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type HttpRequestorM struct {
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewHttpRequestorM() HttpRequestorM {

	r := new(HttpRequestorM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e HttpRequestorM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { // "I_Out"
		e.I_Out(msg, info)
	}
}

func (HttpRequestorM) I_In(msg *messages.SAMessage, info [] *interface{}) {
	// Todo: Transfer HttpRequest message creation to proxy
	inv := msg.Payload.(messages.Invocation)

	request := messages.HttpRequest{
		Method:          inv.Args[0].(string),
		Route:           inv.Op,
		QueryParameters: inv.Args[1].(string),
		Protocol:        "HTTP/1.1",
	}

	// store host & port in 'info'
	*info[0] = inv.Host
	*info[1] = inv.Port

	toCRH := make([]interface{}, 3, 3)
	toCRH[0] = inv.Host
	toCRH[1] = inv.Port
	toCRH[2] = request.Marshal()

	*msg = messages.SAMessage{Payload: toCRH}
}

func (HttpRequestorM) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	//response := messages.HttpResponse{}

	//response := msg.Payload.(messages.HttpResponse)

//	*msg = messages.SAMessage{Payload: response}
}