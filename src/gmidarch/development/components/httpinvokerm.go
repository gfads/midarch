package components

import (
	"apps/httpserver/impl"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
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

	response := messages.HttpResponse{}
	impl.RequestListener(request, &response)

	msgTemp := response.Marshal()
	*msg = messages.SAMessage{Payload: msgTemp}
}
