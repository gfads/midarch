package middleware

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/messages"
	"io"
	"net"
	"shared"
)

//@Type: SRHTCP
//@Behaviour: Behaviour = I_Accept -> P1 ++ P1 = I_Accept -> P1 [] I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> P1
type SRHTCP struct{}

func (s SRHTCP) I_Accept(id string, msg *messages.SAMessage, info *interface{}) {

	infoTemp := *info
	srhInfo := infoTemp.(messages.SRHInfo)

	// check if a listen has been already created
	if srhInfo.Ln == nil { // no listen created
		servAddr, err := net.ResolveTCPAddr("tcp", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		srhInfo.Ln, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	// Accept connections
	conn, err := srhInfo.Ln.Accept()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	srhInfo.Conns = append(srhInfo.Conns, conn)
	srhInfo.CurrentConn = conn

	// Update info
	*info = srhInfo

	// Start goroutine
	go handler(info)

	return
}

func (s SRHTCP) I_Receive(id string, msg *messages.SAMessage, info *interface{}) {

	fmt.Println(shared.GetFunction(),"HERE")

	infoTemp := *info
	infoSrh := infoTemp.(messages.SRHInfo)

	// Receive message from handlers
	tempMsgReceived := <-infoSrh.RcvedMessages
	infoSrh.CurrentConn = tempMsgReceived.Chn

	// Update info
	*info = infoSrh
	msg.Payload = tempMsgReceived.Msg

	return
}

func (s SRHTCP) I_Send(id string, msg *messages.SAMessage, info *interface{}) {
	infoTemp := *info
	srhInfo := infoTemp.(messages.SRHInfo)
	conn := srhInfo.CurrentConn

	msgTemp := msg.Payload.([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo

	return
}

func handler(info *interface{}) {

	infoTemp := *info
	srhInfo := infoTemp.(messages.SRHInfo)
	conn := srhInfo.CurrentConn

	for {
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := conn.Read(size)

		if err == io.EOF {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else if err != nil && err != io.EOF {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = conn.Read(msgTemp)
		if err == io.EOF {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else if err != nil && err != io.EOF {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Chn: conn}

		srhInfo.RcvedMessages <- rcvMessage
	}
}
