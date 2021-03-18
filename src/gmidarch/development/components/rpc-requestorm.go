package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type RPCRequestorM struct {
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewRPCRequestorM() RPCRequestorM {

	r := new(RPCRequestorM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e RPCRequestorM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { // "I_Out"
		e.I_Out(msg, info)
	}
}

func (RPCRequestorM) I_In(msg *messages.SAMessage, info [] *interface{}) {
	inv := msg.Payload.(messages.Invocation)

	//request, err := http.NewRequest(inv.Args[0].(string), inv.Host +":"+ inv.Port + inv.Op +"?"+ inv.Args[1].(string), nil)
	//if err != nil {
	//	fmt.Printf("RPCRequestorM:: %v\n", err)
	//	os.Exit(1)
	//}

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
	toCRH[2] = inv //.Marshal()

	*msg = messages.SAMessage{Payload: toCRH}
}

func (RPCRequestorM) I_Out(msg *messages.SAMessage, info [] *interface{}) {

	//response := messages.HttpResponse{}

	//response := msg.Payload.(*http.Response)

	*msg = messages.SAMessage{Payload: msg.Payload.(int)}
}