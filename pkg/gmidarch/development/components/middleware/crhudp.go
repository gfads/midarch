package middleware

import (
	"reflect"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: CRHUDP
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHUDP struct{}

// func (c CRHUDP) getLocalUDPAddr() *net.UDPAddr {
// 	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
// 	//fmt.Println("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
// 	lib.PrintlnDebug("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
// 	var err error = nil
// 	var localUDPAddr *net.UDPAddr = nil
// 	//shared.LocalAddr = "127.0.0.1:37521"
// 	if shared.LocalAddr != "" {
// 		localUDPAddr, err = net.ResolveUDPAddr("udp", shared.LocalAddr)
// 		if err != nil {
// 			shared.ErrorHandler(shared.GetFunction(), err.Error())
// 		}
// 	}
// 	return localUDPAddr
// }

func (c CRHUDP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
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
	// lib.PrintlnDebug("Will connect", crhInfo.Protocols[addr])
	if _, ok := crhInfo.Protocols[addr]; !ok || reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name() != "UDP" { // no connection open yet
		// lib.PrintlnDebug("Try to connect", crhInfo.Protocols[addr])
		if ok {
			// lib.PrintlnDebug("ElemName", reflect.TypeOf(crhInfo.Protocols[addr]).Elem().Name())
			crhInfo.Protocols[addr].CloseConnection()
		}
		crhInfo.Protocols[addr] = &protocols.UDP{}
		crhInfo.Protocols[addr].ConnectToServer(host, port)

		for {
			time.Sleep(200 * time.Millisecond)
			miopPacket := miop.CreateReqPacket("Connect", []interface{}{shared.AdaptId}, shared.AdaptId) // idx is the Connection ID
			msgPayload := Gobmarshaller{}.Marshall(miopPacket)
			err = crhInfo.Protocols[addr].Send(msgPayload)
			if err != nil {
				lib.PrintlnError("Error on send after dial", crhInfo.Conns[addr], err)
				continue
				//shared.ErrorHandler(shared.GetFunction(), err.Error())
			} //else{
			//	break
			//}
			msgFromServer, err := crhInfo.Protocols[addr].Receive()
			if err != nil {
				// lib.PrintlnDebug("Error while reading Connect msg. Error:", err)
				*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
				crhInfo.Conns[addr].Close()
				crhInfo.Conns[addr] = nil
				delete(crhInfo.Conns, addr)
				return
			}
			if isNewConnection, miopPacket, err := c.isNewConnection(msgFromServer); isNewConnection {
				if miopPacket.Bd.ReqBody.Body[1] == "Ok" {
					break
				}
				if err != nil {
					*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
					crhInfo.Conns[addr].Close()
					crhInfo.Conns[addr] = nil
					delete(crhInfo.Conns, addr)
					return
				}
			}
			//}
		}
	}
	// lib.PrintlnDebug("Connected", crhInfo.Protocols[addr])

	err = crhInfo.Protocols[addr].Send(msgToServer)
	if err != nil {
		lib.PrintlnError("Error trying to send message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		return
	}
	// lib.PrintlnDebug("Sent message", crhInfo.Protocols[addr])

	msgFromServer, err := crhInfo.Protocols[addr].Receive()
	if err != nil {
		lib.PrintlnError("Error trying to read message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		return
	}
	// lib.PrintlnDebug("Received message", crhInfo.Protocols[addr])
	err = VerifyProtocolAdaptation(msgFromServer, crhInfo.Protocols[addr])
	if err != nil {
		lib.PrintlnError("Error verifying adaptation:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Protocols[addr].CloseConnection()
		crhInfo.Protocols[addr] = nil
		delete(crhInfo.Protocols, addr)
		return
	}
	// lib.PrintlnDebug("Adaptation Verified", crhInfo.Protocols[addr])
	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHUDP) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket, error) {
	miop, err := Gobmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop, err
}
