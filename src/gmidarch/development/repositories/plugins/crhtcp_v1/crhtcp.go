package crhtcp

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/components/middleware"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	evolutive "injector"
	"log"
	"net"
	"shared"
)

//@Type: CRHTCP
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHTCP struct {}

func (c CRHTCP) getLocalTcpAddr() (*net.TCPAddr) {
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version 1 adapted")
	fmt.Println("shared.LocalAddr:", shared.LocalAddr)
	log.Println("shared.LocalAddr:", shared.LocalAddr)
	var err error = nil
	var localTCPAddr *net.TCPAddr = nil
	if shared.LocalAddr != "" {
		localTCPAddr, err = net.ResolveTCPAddr("tcp", shared.LocalAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
	return localTCPAddr
}

func (c CRHTCP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version 1 adapted")
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
	fmt.Println("Vai conectar", crhInfo.Conns[addr])
	log.Println("Vai conectar", crhInfo.Conns[addr])
	if _, ok := crhInfo.Conns[addr]; !ok { // no connection open yet
		fmt.Println("Entrou", crhInfo.Conns[addr])
		log.Println("Entrou", crhInfo.Conns[addr])
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}

		localTcpAddr := c.getLocalTcpAddr()
		crhInfo.Conns[addr], err = net.DialTCP("tcp", localTcpAddr, tcpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}
		if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
			fmt.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr())
			log.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
			shared.LocalAddr = crhInfo.Conns[addr].LocalAddr().String()
		}
	}
	fmt.Println("Terminou", crhInfo.Conns[addr])
	log.Println("Terminou", crhInfo.Conns[addr])

	// send message's size
	conn := crhInfo.Conns[addr]
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(),err.Error())
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(),err.Error())
	}

	msgFromServer := c.read(err, conn, size)
	if changeProtocol, miop := c.isAdapt(msgFromServer); changeProtocol {
		log.Println("Adapting, miop.Bd.ReqBody.Body:", miop.Bd.ReqBody.Body)

		evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
		//msgFromServer = c.read(err, conn, size)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHTCP) read(err error, conn net.Conn, size []byte) []byte {
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version 1 adapted")
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
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version 1 adapted")
	miop := middleware.Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}
