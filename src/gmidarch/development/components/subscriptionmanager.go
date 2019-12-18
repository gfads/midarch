package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
)

type Subscriptionmanager struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

type SubscriberRecord struct {
	Host string
	Port string
}

var SubscribersSM map[string][]SubscriberRecord

func NewSubscriptionmanager() Subscriptionmanager {

	r := new(Subscriptionmanager)

	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	return *r
}

func (Subscriptionmanager) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(Subscriptionmanager).I_Process(msg)
}

func (s Subscriptionmanager) I_Process(msg *messages.SAMessage) {
	_req := msg.Payload.(shared.Request)

	switch _req.Op {
	case "Subscribe":
		_topic := _req.Args[0].(string)
		_ip := _req.Args[1].(string)
		_port := _req.Args[2].(string)

		_rTemp := s.Subscribe(_topic, _ip, _port)
		_r := make([]interface{}, 2, 2)
		_r[0] = _rTemp
		_r[1] = SubscribersSM

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "Unsubscribe":
		_topic := _req.Args[0].(string)
		_ip := _req.Args[1].(string)
		_port := _req.Args[2].(string)

		_rTemp := s.Unsubscribe(_topic, _ip, _port)
		_r := make([]interface{}, 2, 2)
		_r[0] = _rTemp
		_r[1] = SubscribersSM

	default:
		fmt.Println("SubscriptionManager:: Operation " + _req.Op + " is not implemented by SubscriptionManager")
		os.Exit(1)
	}
}

func (s Subscriptionmanager) Subscribe(topic string, ip string, port string) bool {
	r := true

	// Check if the list of subscribers has already been created
	if SubscribersSM == nil {
		SubscribersSM = make(map[string][]SubscriberRecord)
	}

	// Check if the topic already exists
	if _, ok := SubscribersSM[topic]; !ok {
		SubscribersSM[topic] = []SubscriberRecord{}
	}

	// Include new subscriber
	SubscribersSM[topic] = append(SubscribersSM[topic], SubscriberRecord{Host: ip, Port: port})

	return r
}

func (s Subscriptionmanager) Unsubscribe(topic string, ip string, port string) bool {
	r := true

	// Check if the list is empty
	if SubscribersSM == nil {
		r = false
	} else {
		records := []SubscriberRecord{}
		ok := false
		if records, ok = SubscribersSM[topic]; !ok {
			r = false
		} else {
			// Remove subscriber
			for i := range records {
				if records[i].Host == ip && records[i].Port == port {
					records[i] = records[len(records)-1] // Replace it with the last one.
					records = records[:len(records)-1]
					SubscribersSM[topic] = records
				}
			}
		}
	}

	return r
}
