package Handlers

import (
	"encoding/json"
	"fmt"
	"gmidarch/development/messages"
	"net"
	"os"
	"shared"
	"strings"
)

type HandlerNotify struct {
	Host string
	Port string
}

func (h HandlerNotify) Start(chn chan interface{}) {
	var conn net.Conn
	var err error
	var ln net.Listener

	// Create server to wait for notifications from 'Notification Consumer'
	addr := shared.ResolveHostIp() + ":" + strings.TrimSpace(h.Port)

	ln, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("HandlerNotify:: Unable to listen at port [%v,%v] \n", h.Host,h.Port)
		os.Exit(1)
	}

	if ln != nil {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Printf("HandlerNotify:: Unable to accept connections at port %v\n", h.Port)
			os.Exit(0)
		}
	}

	// Loop to receive data
	for {
		jsonDecoder := json.NewDecoder(conn)
		msgMOM := messages.MessageMOM{}
		err = jsonDecoder.Decode(&msgMOM)

		if err != nil {
			fmt.Printf("HandlerNotify:: Unable to read data")
			os.Exit(0)
		}
		chn <- msgMOM.Payload
		//fmt.Printf("HandlerNotify:: Received Message :: [%v,%v]\n",conn.LocalAddr(),conn.RemoteAddr())
	}
	return

}

func (h HandlerNotify) StartHandler(chn chan interface{}) {
	go h.Start(chn)
}
