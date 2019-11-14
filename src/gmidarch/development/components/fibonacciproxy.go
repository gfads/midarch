package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Fibonacciproxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

func NewFibonacciproxy() Fibonacciproxy {

	r := new(Fibonacciproxy)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e Fibonacciproxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

/*
var i_PreInvR = make(chan messages.SAMessage)
var i_PosTerR = make(chan messages.SAMessage)

func (e Namingproxy) Fibo(_p1 string, _p2 interface{}) bool {
	_p3 := reflect.ValueOf(_p2).FieldByName("Host").String()
	_p4 := int(reflect.ValueOf(_p2).FieldByName("Port").Int())
	_p5 := reflect.TypeOf(_p2).String()
	_args := []interface{}{_p1, ior.IOR{Host: _p3, Port: _p4, Proxy: _p5, Id: 1313}}
	_reqMsg := messages.SAMessage{message.Invocation{Host: n.Host, Port: n.Port, Op: "Register", Args: _args}}
	i_PreInvR <- _reqMsg

	_repMsg := <-i_PosTerR
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(bool)
	return _reply
}
*/

func (Fibonacciproxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <-i_PreInvR
}

func (Fibonacciproxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerR <- *msg
}

func (Fibonacciproxy) I_InOld(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared.Invocation{Host: "localhost", Port: shared.FIBONACCI_PORT, Req: msg.Payload.(shared.Request)} // TODO

	*msg = messages.SAMessage{Payload: inv}
}

func (Fibonacciproxy) I_OutOld(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: result[0]}
}
