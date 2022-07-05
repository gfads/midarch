package middleware

import (
	"encoding/binary"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	evolutive "injector"
	"log"
	"net"
	"reflect"
	"shared"
)

//@Type: CRHTCP
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.eNot -> Behaviour
type CRHTCP struct {}

func (c CRHTCP) getLocalTcpAddr() (*net.TCPAddr) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	//fmt.Println("shared.LocalAddr:", shared.LocalAddr)
	//log.Println("shared.LocalAddr:", shared.LocalAddr)
	var err error = nil
	var localTCPAddr *net.TCPAddr = nil
	//shared.LocalAddr = "127.0.0.1:37521"
	if shared.LocalAddr != "" {
		localTCPAddr, err = net.ResolveTCPAddr("tcp", shared.LocalAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
	return localTCPAddr
}

func (c CRHTCP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	infoTemp := *info
	crhInfo := infoTemp.(messages.CRHInfo)

	// check message
	//payload := msg.Payload.([]byte)
	payload := msg.Payload.(messages.RequestorInfo).MarshalledMessage
	h := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Host
	p := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Port

	host := ""
	port := ""

	if (h == "" || p == "") {
		host = crhInfo.EndPoint.Host
		port = crhInfo.EndPoint.Port
	} else {
		host = h
		port = p
	}

	msgToServer := payload

	addr := host + ":" + port
	var err error
	//fmt.Println("Vai conectar", crhInfo.Conns[addr])
	//log.Println("Vai conectar", crhInfo.Conns[addr])
	if _, ok := crhInfo.Conns[addr]; !ok || reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name() != "TCPConn" { // no connection open yet
		//fmt.Println("Entrou", crhInfo.Conns[addr])
		//log.Println("Entrou", crhInfo.Conns[addr])
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}

		//localTcpAddr := c.getLocalTcpAddr()
		crhInfo.Conns[addr], err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}
		if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
			//fmt.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr())
			//log.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
			shared.LocalAddr = crhInfo.Conns[addr].LocalAddr().String()
		}
	}
	//fmt.Println("Terminou", crhInfo.Conns[addr])
	//log.Println("Terminou", crhInfo.Conns[addr])

	// send message's size
	conn := crhInfo.Conns[addr]
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	err = c.send(sizeOfMsgSize, msgToServer, conn)
	if err != nil {
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Conns[addr].Close()
		crhInfo.Conns[addr] = nil
		delete(crhInfo.Conns, addr)
		//fmt.Println("Error after trying to send message:", err.Error())
		return
	}

	msgFromServer := c.read(err, conn, sizeOfMsgSize)
	if changeProtocol, miopPacket := c.isAdapt(msgFromServer); changeProtocol {
		//log.Println("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)

		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)

		miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{miopPacket.Bd.ReqBody.Body[0], shared.AdaptId, "Ok"}, shared.AdaptId) // idx is the Connection ID
		msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
		c.send(sizeOfMsgSize, msgPayload, conn)

		if miopPacket.Bd.ReqBody.Body[0] == "udp" {
			log.Println("Adapting => UDP")
			evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tcp" {
			log.Println("Adapting => TCP")
			evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
		} else {
			msgFromServer = c.read(err, conn, sizeOfMsgSize)
			//fmt.Println("=================> ############### ============> ########### TCP: Leu o read")
		}
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHTCP) send(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error {
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := conn.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	return nil
}

func (c CRHTCP) read(err error, conn net.Conn, size []byte) []byte {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	return msgFromServer
}

func (c CRHTCP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}
