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

func (Notificationengine) I_Debug() {
	fmt.Printf("NE:: I_Debug\n")
}

func (e Notificationengine) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	//time.Sleep(1 * time.Second)

	switch op[2] {
	case 'P': // I_Publish
		e.I_Publish(msg, r)
	case 'S': // Subscription Manager
		e.I_SM(msg, r)
	case 'N': // Notification consumer
		e.I_NC(msg, r)
	case 'O': // I_Out
		e.I_Out(msg, r)
	case 'D':
		e.I_Debug()
	}
}

func (e Notificationengine) I_Publish(msg *messages.SAMessage, r *bool) {
	_req := msg.Payload.(shared.Request)

	if _req.Op != "Publish" {
		*r = false
	} else {
		_x1 := _req.Args[0].(messages.MessageMOM)
		//_x2 := _x1["Header"].(map[string]interface{})
		//_topic := _x2["Destination"].(string)
		//_msg := _req.Args[1].(map[string]interface{})
		//_msg := _x2["Payload"].(messages.MessageMOM)
		//fmt.Printf("NE:: %v %v\n",_topic,_msg)

		//fmt.Printf("NE:: I_Publish:: BEGIN:: [%v,%v]\n", _topic,_msg)

		//_r := e.Publish(_topic, _msg)
		_r := e.Publish(_x1.Header.Destination, _x1.Payload)
		//fmt.Printf("NE:: I_Publish:: END:: [%v,%v]\n", _topic,_msg)

		*msg = messages.SAMessage{Payload: _r}
	}
}

func (e Notificationengine) I_NC(msg *messages.SAMessage, r *bool) {

	if len(Topics) == 0 || len(SubscribersNE) == 0 {
		*r = false
	} else {
		currentTopic = nxtTopic()
		_msg := e.Consume(currentTopic)
		if _msg.Msg != nil {
			//fmt.Printf("NE:: I_NC:: BEGIN:: [%v,%v]\n", currentTopic,_msg)
			_args := make([]interface{}, 2, 2)
			_args[0] = _msg
			_args[1] = filterSubscribers(currentTopic)
			*msg = messages.SAMessage{_args}
			//fmt.Printf("NE:: I_NC:: END:: [%v,%v]\n", currentTopic,_msg)
		} else {
			*r = false
		}
	}
}

func (e Notificationengine) I_SM(msg *messages.SAMessage, r *bool) {
	_req := msg.Payload.(shared.Request)

	if _req.Op != "Subscribe" && _req.Op != "Unsubscribe" {
		*r = false
		return
	}
	//fmt.Printf("NE:: I_SM:: %v \n",_req.Op)
}

func (e Notificationengine) I_Out(msg *messages.SAMessage, r *bool) {
	_rTemp := msg.Payload.([]interface{})
	SubscribersNE = _rTemp[1].(map[string][]SubscriberRecord)
	_r := shared.QueueingTermination{R: _rTemp[0].(bool)}

	//fmt.Printf("NE:: I_Out\n")

	*msg = messages.SAMessage{Payload: _r}
}

func (Notificationengine) Consume(topic string) *MessageEnqueued {
	r := new(MessageEnqueued)

	//fmt.Printf("NE:: Consume:: BEGIN:: [%v]\n", topic)
	if _, ok := Topics[topic]; !ok {
		Topics[topic] = make(chan MessageEnqueued, shared.QUEUE_SIZE)
	}

	select {
		case *r = <-Topics[topic]:
			//fmt.Printf("NE:: Consume:: END:: [%v]\n", topic)
			default:
	}
	return r
}

func (Notificationengine) Publish(topic string, msg interface{}) bool {
	r := false

	//fmt.Printf("NE:: Publish:: BEGIN:: [%v,%v]\n", topic,msg)
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

func nxtTopic() string {
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
