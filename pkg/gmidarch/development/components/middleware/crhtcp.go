package middleware

import (
	"fmt"
	"reflect"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: CRHTCP
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHTCP struct{}

// func (c CRHTCP) getLocalTcpAddr() *net.TCPAddr {
// 	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
// 	//fmt.Println("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
// 	lib.PrintlnDebug("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
// 	var err error = nil
// 	var localTCPAddr *net.TCPAddr = nil
// 	//shared.LocalAddr = "127.0.0.1:37521"
// 	if shared.LocalAddr != "" {
// 		localTCPAddr, err = net.ResolveTCPAddr("tcp", shared.LocalAddr)
// 		if err != nil {
// 			shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		}
// 	}
// 	return localTCPAddr
// }

func (c CRHTCP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	infoTemp := *info
	crhInfo := infoTemp.(messages.CRHInfo)

	// check message
	//payload := msg.Payload.([]byte)
	payload := msg.Payload.(messages.RequestorInfo).MarshalledMessage
	h := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Host
	p := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Port

	host := ""
	port := ""

	if h == "" || p == "" {
		host = crhInfo.EndPoint.Host
		port = crhInfo.EndPoint.Port
	} else {
		host = h
		port = p
	}

	msgToServer := payload

	addr := host + ":" + port
	var err error
	//fmt.Println("Will connect", crhInfo.Protocols[addr])
	lib.PrintlnDebug("Will connect", crhInfo.Protocols[addr])
	if _, ok := crhInfo.Protocols[addr]; !ok || reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name() != "TCP" { // no connection open yet
		lib.PrintlnDebug("Try to connect", crhInfo.Protocols[addr])
		if ok {
			lib.PrintlnInfo("ElemName", reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name())
			crhInfo.Protocols[addr].CloseConnection()
		}
		crhInfo.Protocols[addr] = &protocols.TCP{}
		crhInfo.Protocols[addr].ConnectToServer(host, port)
	}
	lib.PrintlnDebug("Connected", crhInfo.Protocols[addr])

	// send message's size

	err = crhInfo.Protocols[addr].Send(msgToServer)
	if err != nil {
		lib.PrintlnError("Error trying to send message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		return
	}
	lib.PrintlnDebug("Sent message", crhInfo.Protocols[addr])

	msgFromServer, err := crhInfo.Protocols[addr].Receive()
	if err != nil {
		lib.PrintlnError("Error trying to read message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		return
	}
	lib.PrintlnDebug("Received message", crhInfo.Protocols[addr])
	VerifyProtocolAdaptation(msgFromServer, crhInfo.Protocols[addr])
	lib.PrintlnDebug("Adaptation Verified", crhInfo.Protocols[addr])
	*msg = messages.SAMessage{Payload: msgFromServer}
}
