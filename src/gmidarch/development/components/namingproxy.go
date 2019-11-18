package components

import (
	"fmt"
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

var i_PreInvRNP = make(chan messages.SAMessage)
var i_PosTerRNP = make(chan messages.SAMessage)

func (e Namingproxy) Register(_p1 string, _p2 interface{}) bool {

	_p3 := reflect.ValueOf(_p2).FieldByName("Host").String()
	_p4 := reflect.ValueOf(_p2).FieldByName("Port").String()
	_p5 := reflect.TypeOf(_p2).String()
	_args := []interface{}{_p1, ior.IOR{Host: _p3, Port: _p4, Proxy: _p5, Id: 1313}}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Register", Args: _args}}

	i_PreInvRNP <- _reqMsg

	_repMsg := <-i_PosTerRNP
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Payload"].(bool)
	return _reply
}
func (e Namingproxy) Lookup(_p1 string) (interface{}, bool) {

	_args := []interface{}{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Lookup", Args: _args}}

	fmt.Printf("Namingproxy:: Lookup :: Here\n")
	i_PreInvRNP <- _reqMsg
	_repMsg := <-i_PosTerRNP

	_payload1 := _repMsg.Payload.(map[string]interface{})
	_payload11 := _payload1["Payload"].([]interface{})

	_ior := _payload11[0].(map[string]interface{})
	_host := _ior["Host"].(string)
	_port := _ior["Port"].(string)
	_proxy := _ior["Proxy"].(string)

	_r2 := _payload11[1].(bool)

	if _r2 { // service is registered in naming service
		_p := ProxyLibrary[_proxy]
		proxyPointer := reflect.New(_p)
		proxyValue := proxyPointer.Elem()
		proxyValue.FieldByName("Host").SetString(_host)
		proxyValue.FieldByName("Port").SetString(_port)
		_r1 := proxyValue.Interface()
		return _r1, _r2
	} else { // service NOT registered in naming service
		return nil, _r2
	}
}
func (e Namingproxy) List() [] interface{} {

	_args := []interface{}{}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "List", Args: _args}}

	i_PreInvRNP <- _reqMsg

	_repMsg := <-i_PosTerRNP

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Payload"].([]interface{})
	return _reply
}

func (Namingproxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <-i_PreInvRNP
}
func (Namingproxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerRNP <- *msg
}
