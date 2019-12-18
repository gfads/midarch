package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Notificationengine struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

type MessageEnqueued struct {
	Msg interface{}
}

var SubscribersNE = map[string][]SubscriberRecord{}
var Topics = map[string]chan MessageEnqueued{}
var MsgToBeNotified string
var TopicToBePublished string

func Newnotificationengine() Notificationengine {

	r := new(Notificationengine)
	r.Behaviour = "B = InvP.e1 -> (I_SM -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B [] I_NC -> InvR.e3 -> TerR.e3 -> TerP.e1 -> B)"

	return *r
}

func (e Notificationengine) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'S' { // Subscription Manager
		e.I_SM(msg, r)
		return
	} else if op[2] == 'N' { // Notification consumer
		e.I_NC(msg, r)
		return
	} else if op[2] == 'O' { // I_Out
		e.I_Out(msg, r)
	}
}

func (e Notificationengine) I_Out(msg *messages.SAMessage, r *bool) {
	_rTemp := msg.Payload.(shared.QueueingTermination)
	_rTemp1 := _rTemp.R.([]interface{})
	SubscribersNE = _rTemp1[1].(map[string][]SubscriberRecord)
	_r := shared.QueueingTermination{R: _rTemp1[0].(bool)}

	*msg = messages.SAMessage{Payload: _r}
}

func (e Notificationengine) I_SM(msg *messages.SAMessage, r *bool) {
	request := msg.Payload.(shared.Request)

	if request.Op != "Subscribe" && request.Op != "Unsubscribe" {
		*r = false
	}
}

func (e Notificationengine) I_NC(msg *messages.SAMessage, r *bool) {
	request := msg.Payload.(shared.Request)

	if request.Op != "Publish" {
		*r = false
	} else {
		// send 'Subscribers' to Notification Consumer
		_args := make([]interface{}, 2, 2)
		_args[0] = *msg
		_args[1] = SubscribersNE

		*msg = messages.SAMessage{Payload: _args}
	}
}