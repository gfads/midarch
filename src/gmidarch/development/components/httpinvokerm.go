package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"strings"
)

type HttpInvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewHttpInvokerM() HttpInvokerM {
	r := new(HttpInvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e HttpInvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
}

func (HttpInvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) { // TODO
	// unmarshall
	payload := msg.Payload.([]byte)
	// Todo: Remove manual unmarshall HTTP message from here
	message := messages.HttpMessage{}
	message.Headers.Fields = make(map[string]string)
	lines := strings.Split(string(payload), "\n")
	bodyStarted := false
	for _, line := range lines {
		if message.Method == "" {
			startLine := strings.Fields(line)
			message.Method = startLine[0]
			message.Protocol = startLine[2]

			route := strings.Split(startLine[1], "?")
			message.Route = route[0]
			if len(route) > 1 {
				message.QueryParameters = route[1]
			}

			continue
		}

		if strings.TrimSpace(line) == "" {
			bodyStarted = true
		}

		if !bodyStarted {
			header := strings.Split(line, ": ")
			fmt.Println("HttpInvokerM.I_Process header:", header)
			message.Headers.Fields[header[0]] = header[1]
		}else{
			message.Body += line
		}
	}

	fmt.Println("HttpInvokerM.I_Process method:", message.Method)
	switch message.Method {
	case "GET":
		msgTemp := []byte(`HTTP/1.1 200 OK
content-type: text/html; charset=UTF-8
date: Sun 06 Sep 2020 14:39:08 GMT

<html><h1>Teste ok</h1></html>`)
		fmt.Println("HttpInvokerM.I_Process GET:", msgTemp)
		*msg = messages.SAMessage{Payload: msgTemp}
	default:
		fmt.Printf("HttpInvokerM:: Operation '%v' not implemented by Http Service\n", message.Method)
		os.Exit(0)
	}
	fmt.Println("HttpInvokerM.I_Process finished")
}
