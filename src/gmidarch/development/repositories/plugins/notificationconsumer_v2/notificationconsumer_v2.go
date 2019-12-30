package main

import (
	"encoding/json"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/components"
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
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1  -> B"

	return *r
}

func Gettype() interface{} {
	return Notificationconsumer{}
}

func Getselector() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool){
	return Notificationconsumer{}.Selector
}

func (e Notificationconsumer) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg)
}

var ActiveConsumers = map[string]bool{}

var connNC map[string]net.Conn

func (n Notificationconsumer) I_Process(msg *messages.SAMessage) {
	_args := msg.Payload.([]interface{})
	_msg := _args[0].(*components.MessageEnqueued)
	_subscribers := _args[1].([]components.SubscriberRecord)

	for s := range _subscribers {
		n.NotifySubscriber(_subscribers[s], *_msg)
	}
	_r := true // TODO, check if all subscribers were actually notified
	*msg = messages.SAMessage{_r}
}

func (Notificationconsumer) NotifySubscriber(subscriber components.SubscriberRecord, msg components.MessageEnqueued) {

	err := *new(error)
	// Check if 'Active Consumers' (Consumers whose connection to Handler already exists) has been created
	//if ActiveConsumers == nil {
	//	ActiveConsumers = make(map[string]bool, shared.MAX_NUMBER_OF_ACTIVE_CONSUMERS)
	//}
	if connNC == nil {
		connNC = make(map[string] net.Conn, shared.MAX_NUMBER_OF_ACTIVE_CONSUMERS)
	}

	// Notify Subscribers
	host := subscriber.Host
	port := subscriber.Port
	addr := strings.Join([]string{host, port}, ":")

	// Check if the connection with the Handler already exists
	//fmt.Printf("NC:: ActiveConsumer:: [%v,%v]\n",host,port)
	//_, ok := ActiveConsumers[addr]
	_, ok := connNC[addr]
	if !ok {
		//ActiveConsumers[addr] = true
		connNC[addr], err = net.Dial("tcp", addr)

		//portTmp := port  // TODO
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
	//fmt.Printf("NC:: message sent to [%v,%v,%v]\n",connNC[addr].LocalAddr(),connNC[addr].RemoteAddr(),msgMOM)
}