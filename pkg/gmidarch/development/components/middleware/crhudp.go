package middleware

import (
	"encoding/binary"
	"net"
	"reflect"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: CRHUDP
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHUDP struct{}

func (c CRHUDP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	infoTemp := *info
	crhInfo := infoTemp.(messages.CRHInfo)
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)

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
	//var err error
	if _, ok := crhInfo.Conns[addr]; !ok || reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name() != "UDPConn" { // no connection open yet
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		localUdpAddr := c.getLocalUdpAddr()

		crhInfo.Conns[addr], err = net.DialUDP("udp", localUdpAddr, udpAddr)
		if err != nil {
			lib.PrintlnError("Dial error", crhInfo.Conns[addr], err)
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		} //else{

		for {
			time.Sleep(200 * time.Millisecond)
			miopPacket := miop.CreateReqPacket("Connect", []interface{}{shared.AdaptId}, shared.AdaptId) // idx is the Connection ID
			msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
			err = c.send(sizeOfMsgSize, msgPayload, crhInfo.Conns[addr])
			if err != nil {
				lib.PrintlnError("Error on send after dial", crhInfo.Conns[addr], err)
				continue
				//shared.ErrorHandler(shared.GetFunction(), err.Error())
			} //else{
			//	break
			//}
			msgFromServer, err := c.read(crhInfo.Conns[addr], sizeOfMsgSize)
			if err != nil {
				lib.PrintlnDebug("Error while reading Connect msg. Error:", err)
				*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
				crhInfo.Conns[addr].Close()
				crhInfo.Conns[addr] = nil
				delete(crhInfo.Conns, addr)
				return
			}
			if isNewConnection, miopPacket := c.isNewConnection(msgFromServer); isNewConnection {
				if miopPacket.Bd.ReqBody.Body[1] == "Ok" {
					break
				}
			}
			//}
		}

		if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
			//fmt.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr())
			lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
			shared.LocalAddr = crhInfo.Conns[addr].LocalAddr().String()
		}
	}

	// send message's size
	conn := crhInfo.Conns[addr]
	c.send(sizeOfMsgSize, msgToServer, conn)

	msgFromServer, err := c.read(conn, sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError("Error trying to read message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Conns[addr].Close()
		crhInfo.Conns[addr] = nil
		delete(crhInfo.Conns, addr)
		return
	}
	if changeProtocol, miopPacket := c.isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)
		//lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body[0]:", miopPacket.Bd.ReqBody.Body[0])
		//lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body[1]:", miopPacket.Bd.ReqBody.Body[1])
		//lib.PrintlnDebug("Adapting, shared.AdaptId:", shared.AdaptId)
		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)

		miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{miopPacket.Bd.ReqBody.Body[0], shared.AdaptId, "Ok"}, shared.AdaptId) // idx is the Connection ID
		msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
		c.send(sizeOfMsgSize, msgPayload, conn)

		if miopPacket.Bd.ReqBody.Body[0] == "udp" {
			lib.PrintlnInfo("Adapting => UDP")
			//evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhudp")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tcp" {
			lib.PrintlnInfo("Adapting => TCP")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtcp")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tls" {
			lib.PrintlnInfo("Adapting => TLS")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhtls")
		} else if miopPacket.Bd.ReqBody.Body[0] == "quic" {
			lib.PrintlnInfo("Adapting => QUIC")
			//evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, "crhquic")
		} else {
			msgFromServer, _ = c.read(conn, sizeOfMsgSize)
			//fmt.Println("=================> ############### ============> ########### TCP: Leu o read")
		}
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Leu")

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHUDP) send(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error {
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := conn.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		lib.PrintlnError("Erro no send 1, retornou o erro", err)
		return err
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Escreveu sizeOfMsgSize")

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		//fmt.Println("Erro no envio do sizeOfMsgSize(", sizeOfMsgSize, ") Connection:", reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name())
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		lib.PrintlnError("Erro no send 2, retornou o erro", err)
		return err
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Escreveu msg")
	return nil
}

func (c CRHUDP) getLocalUdpAddr() *net.UDPAddr {
	var err error = nil
	var localUdpAddr *net.UDPAddr = nil
	//shared.LocalAddr = "127.0.0.1:37522"
	if shared.LocalAddr != "" {
		//fmt.Println("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
		lib.PrintlnDebug("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
		localUdpAddr, err = net.ResolveUDPAddr("udp", shared.LocalAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	} //else{
	//fmt.Println("else shared.LocalAddr:", shared.LocalAddr)
	//}
	return localUdpAddr
}

func (c CRHUDP) read(conn net.Conn, size []byte) ([]byte, error) {
	// receive reply's size
	err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
	}
	_, err = conn.Read(size)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return nil, err
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	err = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
	}
	_, err = conn.Read(msgFromServer)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return nil, err
	}
	return msgFromServer, nil
}

func (c CRHUDP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (c CRHUDP) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
