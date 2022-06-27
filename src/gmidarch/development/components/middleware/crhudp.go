package middleware

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	evolutive "injector"
	"log"
	"net"
	"reflect"
	"shared"
)

//@Type: CRHUDP
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHUDP struct {}

func (c CRHUDP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
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
	if _, ok := crhInfo.Conns[addr]; !ok || reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name() != "UDPConn" { // no connection open yet
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}

		localUdpAddr := c.getLocalUdpAddr()
		crhInfo.Conns[addr], err = net.DialUDP("udp", localUdpAddr, udpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
			fmt.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr())
			log.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
			shared.LocalAddr = crhInfo.Conns[addr].LocalAddr().String()
		}
	}

	// send message's size
	conn := crhInfo.Conns[addr]
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Escreveu size")

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Println("Erro no envio do size(", size, ") Connection:", reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name())
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Escreveu msg")

	msgFromServer := c.read(err, conn, size)
	if changeProtocol, miop := c.isAdapt(msgFromServer); changeProtocol {
		log.Println("Adapting, miop.Bd.ReqBody.Body:", miop.Bd.ReqBody.Body)
		log.Println("Adapting, miop.Bd.ReqBody.Body[0]:", miop.Bd.ReqBody.Body[0])
		log.Println("Adapting, miop.Bd.ReqBody.Body[1]:", miop.Bd.ReqBody.Body[1])
		log.Println("Adapting, shared.AdaptId:", shared.AdaptId)
		shared.AdaptId = miop.Bd.ReqBody.Body[1].(int)
		if miop.Bd.ReqBody.Body[0] == "udp" {
			log.Println("Adapting => UDP")
			evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
		} else if miop.Bd.ReqBody.Body[0] == "tcp" {
			log.Println("Adapting => TCP")
			evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
		} else {
			msgFromServer = c.read(err, conn, size)
		}
	}
	log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted ###### Leu")

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHUDP) getLocalUdpAddr() (*net.UDPAddr) {
	var err error = nil
	var localUdpAddr *net.UDPAddr = nil
	if shared.LocalAddr != "" {
		fmt.Println("shared.LocalAddr:", shared.LocalAddr)
		log.Println("shared.LocalAddr:", shared.LocalAddr)
		localUdpAddr, err = net.ResolveUDPAddr("udp", shared.LocalAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
	return localUdpAddr
}

func (c CRHUDP) read(err error, conn net.Conn, size []byte) []byte {
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

func (c CRHUDP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}
