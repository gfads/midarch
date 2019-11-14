package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"reflect"
	"shared/ior"
)

type Namingproxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

func NewNamingproxy() Namingproxy {

	r := new(Namingproxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e Namingproxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

var i_PreInvR = make(chan messages.SAMessage)
var i_PosTerR = make(chan messages.SAMessage)

func (e Namingproxy) Register(_p1 string, _p2 interface{}) bool {

	_p3 := reflect.ValueOf(_p2).FieldByName("Host").String()
	_p4 := reflect.ValueOf(_p2).FieldByName("Port").String()
	_p5 := reflect.TypeOf(_p2).String()
	_args := []interface{}{_p1, ior.IOR{Host: _p3, Port: _p4, Proxy: _p5, Id: 1313}}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Register", Args: _args}}

	i_PreInvR <- _reqMsg

	_repMsg := <-i_PosTerR
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(bool)
	return _reply
}

func (Namingproxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <-i_PreInvR

	//fmt.Printf("NamingProxy:: HERE\n")
}

func (Namingproxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerR <- *msg
}
