package components

import (
	"encoding/json"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"net"
	"os"
	"shared"
	"strings"
)

type NotificationengineX struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

type MessageEnqueued struct {
	Msg interface{}
}

type SubscriberRecord struct {
	Host string
	Port string
}

const MAX_NUMBER_OF_TOPICS = 3

var connNC map[string]net.Conn
var Subscribers = map[string][]SubscriberRecord{}
var Topics = map[string]chan MessageEnqueued{}
var currentTopic = ""

func NewnotificationengineX() NotificationengineX {

	r := new(NotificationengineX)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B [] I_NC -> B"
	return *r
}

func (e NotificationengineX) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'P' {
		e.I_Process(msg)
	} else {
		e.I_NC(msg, r)
	}
}

func (e NotificationengineX) I_Process(msg *messages.SAMessage) {
	_req := msg.Payload.(shared.Request)

	switch _req.Op {
	case "Publish":
		_momMessage := _req.Args[0].(messages.MessageMOM)
		_r := e.Publish(_momMessage.Header.Destination, _momMessage.Payload) // Publish
		*msg = messages.SAMessage{Payload: _r}
	case "Subscribe":
		_topic := _req.Args[0].(string)
		_ip := _req.Args[1].(string)
		_port := _req.Args[2].(string)
		_r := e.Subscribe(_topic, _ip, _port) // Subscribe
		*msg = messages.SAMessage{_r}
	case "Unsubscribe":
		_topic := _req.Args[0].(string)
		_ip := _req.Args[1].(string)
		_port := _req.Args[2].(string)
		_r := e.Unsubscribe(_topic, _ip, _port) // Unsubscribe
		*msg = messages.SAMessage{_r}
	default:
		fmt.Println("Notificationengine:: Operation " + _req.Op + " is not implemented by Notificationengine")
		os.Exit(1)
	}
}

func (e NotificationengineX) I_NC(msg *messages.SAMessage, r *bool) {

	if len(Topics) == 0 || len(Subscribers) == 0 {
		*r = false
		return
	} else {
		currentTopic = nextTopic()
		_msg := e.Consume(currentTopic)
		if _msg.Msg != nil {
			_subscribers := filterSubscribers(currentTopic)
			for s := range _subscribers {
				e.NotifySubscriber(_subscribers[s], *_msg)
			}
		} else {
			*r = false
			return
		}
	}
}

func (s NotificationengineX) Subscribe(topic string, ip string, port string) bool {
	r := true

	// Check if the list of subscribers has already been created
	if Subscribers == nil {
		Subscribers = make(map[string][]SubscriberRecord)
	}

	// Check if the topic already exists
	if _, ok := Subscribers[topic]; !ok {
		Subscribers[topic] = []SubscriberRecord{}
	}

	// Include new subscriber
	Subscribers[topic] = append(Subscribers[topic], SubscriberRecord{Host: ip, Port: port})

	return r
}

func (s NotificationengineX) Unsubscribe(topic string, ip string, port string) bool {
	r := true

	// Check if the list is empty
	if Subscribers == nil {
		r = false
	} else {
		records := []SubscriberRecord{}
		ok := false
		if records, ok = Subscribers[topic]; !ok {
			r = false
		} else {
			// Remove subscriber
			for i := range records {
				if records[i].Host == ip && records[i].Port == port {
					records[i] = records[len(records)-1] // Replace it with the last one.
					records = records[:len(records)-1]
					Subscribers[topic] = records
				}
			}
		}
	}

	return r
}

func (NotificationengineX) Consume(topic string) *MessageEnqueued {
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

func (NotificationengineX) Publish(topic string, msg interface{}) bool {
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

	for t := range Subscribers {
		if t == topic {
			r = Subscribers[t]
			break
		}
	}
	return r
}

func (NotificationengineX) NotifySubscriber(subscriber SubscriberRecord, msg MessageEnqueued) {

	err := *new(error)
	if connNC == nil {
		connNC = make(map[string]net.Conn, shared.MAX_NUMBER_OF_ACTIVE_CONSUMERS)
	}

	// Notify Subscribers
	host := subscriber.Host
	port := subscriber.Port
	addr := strings.Join([]string{host, port}, ":")

	// Check if the connection with the Handler already exists
	_, ok := connNC[addr]
	if !ok {
		connNC[addr], err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("NofiticationConsumer:: %v\n", err)
			os.Exit(0)
		}
	}

	// Prepare message to be sent to Handler
	msgMOM := messages.MessageMOM{Header: messages.MOMHeader{""}, Payload: msg}

	// Send message
	encoder := json.NewEncoder(connNC[addr]) // TODO
	err = encoder.Encode(msgMOM)
	if err != nil {
		fmt.Printf("NofiticationConsumer:: %v\n", err)
		os.Exit(0)
	}
}
