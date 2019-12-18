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

type Notificationconsumer struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newnotificationconsumer() Notificationconsumer {

	r := new(Notificationconsumer)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e Notificationconsumer) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg)
}

var ActiveConsumers = map[string]bool{}

var connNC net.Conn

func (n Notificationconsumer) I_Process(msg *messages.SAMessage) {
	_args := msg.Payload.([]interface{})
	_req := _args[0].(messages.SAMessage).Payload.(shared.Request)
	_subscribers := _args[1].(map[string][]SubscriberRecord)
	_topic := _req.Args[0].(string)
	_msg := _req.Args[1].(string) // TODO

	switch _req.Op {
	case "Publish":
		n.NotifySubscribers(_subscribers,_topic,_msg)
		_r := true // TODO, check if all subscribers were actually notified

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	default:
		fmt.Println("NotificationConsumer:: Operation " + _req.Op + " is not implemented by NotificationConsumer")
		os.Exit(0)
	}
}

func (Notificationconsumer) NotifySubscribers(subscribers map[string][]SubscriberRecord, topic string,msg string) {

	err := *new(error)
	// Check if 'Active Consumers' (Consumers whose connection to Handler already exists) has been created
	if ActiveConsumers == nil {
		ActiveConsumers = make(map[string]bool, shared.MAX_NUMBER_OF_ACTIVE_CONSUMERS)
	}

	// Notify Subscribers
	filteredSubscribers := filterSubscribers(subscribers, topic)
	for i := range filteredSubscribers {
		host := filteredSubscribers[i].Host
		port := filteredSubscribers[i].Port
		addr := strings.Join([]string{host, port}, ":")

		// Check if the connection with the Handler already exists
		_, ok := ActiveConsumers[addr]
		if !ok {
			ActiveConsumers[addr] = true
			connNC, err = net.Dial("tcp", addr)

			//portTmp := port  // TODO
			if err != nil {
				fmt.Printf("NofiticationConsumer:: %v\n",err)
				os.Exit(0)
			}
		}

		// Prepare message to be sent to Handler
		msgMOM := messages.MessageMOM{Header: messages.MOMHeader{""}, Payload: msg}

		// Send message
		encoder := json.NewEncoder(connNC) // TODO
		err = encoder.Encode(msgMOM)
		if err != nil {
			fmt.Printf("NofiticationConsumer:: %v\n",err)
			os.Exit(0)
		}
	}

	return
}

func filterSubscribers(subscribers map[string][]SubscriberRecord, topic string) []SubscriberRecord {
	r := []SubscriberRecord{}

	for t := range subscribers {
		if t == topic {
			r = subscribers[t]
			break
		}
	}
	return r
}