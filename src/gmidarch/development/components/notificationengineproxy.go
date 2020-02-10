package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
	"shared/Handlers"
)

type Notificationengineproxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

var i_PreInvNEP = make(chan messages.SAMessage)
var i_PosTerNEP = make(chan messages.SAMessage)
var HandlersProxy = make(map[string]Handlers.HandlerNotify,10)  // TODO

func Newnotificationengineproxy() Notificationengineproxy {

	r := new(Notificationengineproxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e Notificationengineproxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg)
	} else { //"I_Out"
		e.I_Out(msg)
	}
}

func (e Notificationengineproxy) Publish(_p1 string, _p2 interface{}) bool {
	_header := messages.MOMHeader{Destination:_p1}
	_msgMOM := messages.MessageMOM{Header:_header,Payload:_p2}
	_args := []interface{}{_msgMOM}

	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Publish", Args: _args}}

	i_PreInvNEP  <- _reqMsg

	//fmt.Printf("NotificationProxy:: Publish:: Before response \n")

	_repMsg := <-i_PosTerNEP

	//fmt.Printf("NotificationProxy:: Publish:: After response\n")

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Payload"].(bool)

	return _reply
}

func (e Notificationengineproxy) Subscribe(_p1 string, _chn chan interface{}) (bool) {
	_p2 := shared.ResolveHostIp()             // host TODO
	_p3 := shared.NextPortTCPAvailable()      // port
	_args := []interface{}{_p1,_p2,_p3}

	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Subscribe", Args: _args}}

	i_PreInvNEP  <- _reqMsg
	_repMsg := <-i_PosTerNEP

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Payload"].(bool)

	// Create the new handler associated to the topic [one handler per topic]
	if _,ok := HandlersProxy[_p1]; !ok{
		fmt.Printf("NE:: %v\n",_p2)
		HandlersProxy[_p1] = Handlers.HandlerNotify{Host:_p2,Port:_p3}
	}

	// Start the new handler
	HandlersProxy[_p1].StartHandler(_chn)

	return _reply
}

func (Notificationengineproxy) I_In(msg *messages.SAMessage) {
	*msg = <- i_PreInvNEP
}

func (Notificationengineproxy) I_Out(msg *messages.SAMessage) {
	i_PosTerNEP <- *msg
}