package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"net/http"
	"os"
)

type Http2RequestorM struct {
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewHttp2RequestorM() Http2RequestorM {

	r := new(Http2RequestorM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e Http2RequestorM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { // "I_Out"
		e.I_Out(msg, info)
	}
}

func (Http2RequestorM) I_In(msg *messages.SAMessage, info [] *interface{}) {
	// Todo: Transfer HttpRequest message creation to proxy
	inv := msg.Payload.(messages.Invocation)

	request, err := http.NewRequest(inv.Args[0].(string), inv.Host +":"+ inv.Port + inv.Op +"?"+ inv.Args[1].(string), nil)
	if err != nil {
		fmt.Printf("Http2RequestorM:: %v\n", err)
		os.Exit(1)
	}

	//		messages.HttpRequest{
	//	Method:          inv.Args[0].(string),
	//	Route:           inv.Op,
	//	QueryParameters: inv.Args[1].(string),
	//	Protocol:        "HTTP/1.1",
	//}

	// store host & port in 'info'
	*info[0] = inv.Host
	*info[1] = inv.Port

	toCRH := make([]interface{}, 3, 3)
	toCRH[0] = inv.Host
	toCRH[1] = inv.Port
	toCRH[2] = request //.Marshal()

	*msg = messages.SAMessage{Payload: toCRH}
}

func (Http2RequestorM) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	// Todo: Implement HTTPS
	//response := messages.HttpResponse{}

	response := msg.Payload.(*http.Response)

	*msg = messages.SAMessage{Payload: response}
}