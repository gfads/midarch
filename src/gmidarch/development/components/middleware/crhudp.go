package middleware

import (
	"encoding/binary"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	evolutive "injector"
	"log"
	"net"
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

	key := host + ":" + port
	var err error
	if _, ok := crhInfo.Conns[key]; !ok { // no connection open yet
		udpAddr, err := net.ResolveUDPAddr("udp", key)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}

		crhInfo.Conns[key], err = net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}
	}

	// send message's size
	conn := crhInfo.Conns[key]
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
		if miop.Bd.ReqBody.Body[0] == "udp" {
			evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
		}
		msgFromServer = c.read(err, conn, size)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
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
