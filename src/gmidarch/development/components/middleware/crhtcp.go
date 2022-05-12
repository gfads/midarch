package middleware

import (
	"encoding/binary"
	"gmidarch/development/messages"
	"net"
	"shared"
)

//@Type: CRHTCP
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHTCP struct {}

func (c CRHTCP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
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
		tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(),err.Error())
		}

		crhInfo.Conns[key], err = net.DialTCP("tcp", nil, tcpAddr)
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

	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(),err.Error())
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(),err.Error())
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}


