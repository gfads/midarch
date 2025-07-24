package middleware

import (
	"reflect"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: CRHRPC
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHRPC struct{}

// func (c CRHRPC) getLocalTcpAddr() *net.TCPAddr {
// 	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
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

func (c CRHRPC) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
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
	lib.PrintlnDebug("Will connect", crhInfo.Conns[addr])
	if _, ok := crhInfo.Protocols[addr]; !ok || reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name() != "RPC" { // no connection open yet
		lib.PrintlnDebug("Try to connect", crhInfo.Protocols[addr])
		if ok {
			lib.PrintlnDebug("ElemName", reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name())
			crhInfo.Protocols[addr].CloseConnection()
		}
		crhInfo.Protocols[addr] = &protocols.RPC{}
		crhInfo.Protocols[addr].ConnectToServer(host, port)
	}
	lib.PrintlnDebug("Connected", crhInfo.Protocols[addr])

	err = crhInfo.Protocols[addr].Send(msgToServer)
	if err != nil {
		lib.PrintlnError("Error trying to send message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		// lib.SHOW_MESSAGES = append(lib.SHOW_MESSAGES, lib.DEBUG)
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

// func (c CRHRPC) send(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error {
// 	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
// 	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
// 	_, err := conn.Write(sizeOfMsgSize)
// 	if err != nil {
// 		//shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		return err
// 	}

// 	// send message
// 	_, err = conn.Write(msgToServer)
// 	if err != nil {
// 		//shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		return err
// 	}
// 	return nil
// }

// func (c CRHRPC) read(conn net.Conn, size []byte) ([]byte, error) {
// 	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHRPC Version Not adapted")
// 	// receive reply's size
// 	_, err := conn.Read(size)
// 	if err != nil {
// 		lib.PrintlnError(shared.GetFunction(), err)
// 		//shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		return nil, err
// 	}

// 	// receive reply
// 	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
// 	_, err = conn.Read(msgFromServer)
// 	if err != nil {
// 		lib.PrintlnError(shared.GetFunction(), err)
// 		//shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		return nil, err
// 	}
// 	return msgFromServer, nil
// }
