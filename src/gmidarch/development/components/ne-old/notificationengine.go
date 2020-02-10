package components

import (
	"fmt"
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

const MAX_NUMBER_OF_TOPICS = 3

var SubscribersNE = map[string][]SubscriberRecord{}
var Topics = map[string]chan MessageEnqueued{}
var currentTopic = ""

func Newnotificationengine() Notificationengine {

	r := new(Notificationengine)
	r.Behaviour = "B = I_NC -> InvR.e2 -> TerR.e2 -> B [] InvP.e1 -> (I_SM -> InvR.e3 -> TerR.e3 -> I_Out -> TerP.e1 -> B [] I_Publish -> TerP.e1 -> B)"
	return *r
}

func (e Notificationengine) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	//time.Sleep(10 * time.Millisecond)

	switch op[2] {
	case 'P': // I_Publish
		e.I_Publish(msg, r)
	case 'S': // Subscription Manager
		e.I_SM(msg, r)
	case 'N': // Notification consumer
		e.I_NC(msg, r)
	case 'O': // I_Out
		e.I_Out(msg, r)
	}
}

func (e Notificationengine) I_Publish(msg *messages.SAMessage, r *bool) {
	_req := msg.Payload.(shared.Request)

	if _req.Op != "Publish" {
		*r = false
	} else {
		_x1 := _req.Args[0].(messages.MessageMOM)
		_r := e.Publish(_x1.Header.Destination, _x1.Payload)

		*msg = messages.SAMessage{Payload: _r}
	}
}

func (e Notificationengine) I_NC(msg *messages.SAMessage, r *bool) {

	fmt.Printf("NE:: Current Topic:: [%v] Connection:: [%v]\n", currentTopic, currentConnection)
	if len(Topics) == 0 || len(SubscribersNE) == 0 {
		*r = false
		return
	} else {
		currentTopic = nextTopic()
		_msg := e.Consume(currentTopic)
		if _msg.Msg != nil {
			_args := make([]interface{}, 2, 2)
			_args[0] = _msg
			_args[1] = filterSubscribers(currentTopic)
			*msg = messages.SAMessage{_args}
		} else {
			*r = false
			return
		}
	}
}

func (e Notificationengine) I_SM(msg *messages.SAMessage, r *bool) {
	_req := msg.Payload.(shared.Request)

	if _req.Op != "Subscribe" && _req.Op != "Unsubscribe" {
		*r = false
		return
	}
}

func (e Notificationengine) I_Out(msg *messages.SAMessage, r *bool) {
	_rTemp := msg.Payload.([]interface{})
	SubscribersNE = _rTemp[1].(map[string][]SubscriberRecord)
	_r := shared.QueueingTermination{R: _rTemp[0].(bool)}

	*msg = messages.SAMessage{Payload: _r}
}

func (Notificationengine) Consume(topic string) *MessageEnqueued {
	r := new(MessageEnqueued)

	if _, ok := Topics[topic]; !ok {
		Topics[topic] = make(chan MessageEnqueued, shared.QUEUE_SIZE)
	}

	select {
	case *r = <-Topics[topic]:
	default:
	}
	return r
}

func (Notificationengine) Publish(topic string, msg interface{}) bool {
	r := false

	// Check if the topic exist
	if _, ok := Topics[topic]; !ok {
		Topics[topic] = make(chan MessageEnqueued, shared.QUEUE_SIZE)
	}

	// Put the message on the topic
	if len(Topics[topic]) < shared.QUEUE_SIZE {
		Topics[topic] <- MessageEnqueued{Msg: msg}
		r = true
	} else {
		r = false
	}
	//fmt.Printf("NE:: Publish:: END:: [%v,%v]\n", topic,msg)
	return r
}

func nextTopic() string {
	r := ""

	if currentTopic == "" {
		if len(Topics) != 0 {
			for i := range Topics {
				r = i // return first topic in the Topics
				break
			}
		}
	} else { // need to discover the next topic in the current map
		// find current keys
		keys := make([]string, len(Topics))
		idx := 0
		for i := range Topics {
			keys[idx] = i
			idx++
		}
		// find next key
		for i := 0; i < len(keys); i++ {
			if keys[i] == currentTopic {
				if i+1 < len(Topics) {
					r = keys[i+1]
					break
				} else {
					r = keys[0]
					break
				}
			}
		}
	}
	return r
}

func filterSubscribers(topic string) []SubscriberRecord {
	r := []SubscriberRecord{}

	for t := range SubscribersNE {
		if t == topic {
			r = SubscribersNE[t]
			break
		}
	}
	return r
}
