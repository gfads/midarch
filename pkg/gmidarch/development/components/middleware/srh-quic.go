package middleware

import (
	"io"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"github.com/quic-go/quic-go"
)

// @Type: SRHQUIC
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHQUIC struct{}

func (s SRHQUIC) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHQUIC Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	if srhInfo.Protocol == nil {
		srhInfo.Protocol = &protocols.QUIC{}
		srhInfo.Protocol.StartServer(srhInfo.EndPoint.Host, srhInfo.EndPoint.Port, 2) //shared.MAX_NUMBER_OF_CONNECTIONS)
	}

	// // check if a listener has already been created
	// if srhInfo.Ln == nil { // no listen created
	// 	servAddr, err := net.ResolveTCPAddr("tcp", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
	// 	if err != nil {
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// 	srhInfo.Ln, err = net.ListenTLS("tcp", servAddr)
	// 	if err != nil {
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// }

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := srhInfo.Protocol.AvailableConnectionFromPool()
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHQUIC Version 2 adapted - No connection available")
		time.Sleep(1 * time.Millisecond)
		return
	}

	go func() {
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

		client := srhInfo.Protocol.WaitForConnection(availableConenctionIndex)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

		// Update info
		*info = srhInfo

		// Start goroutine
		if client != nil {
			go s.handler(info, availableConenctionIndex)
		}
	}()
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHQUIC Version Not adapted")
	return
}

func (s SRHQUIC) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHQUIC Version Not adapted")
	//lib.PrintlnDebug(shared.GetFunction(), "HERE")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)

	select {
	case tempMsgReceived := <-srhInfo.RcvedMessages:
		{
			// Receive message from handlers
			//srhInfo.CurrentConn = tempMsgReceived.Conn

			// Update info
			*info = srhInfo
			msg.Payload = tempMsgReceived.Msg
			//lib.PrintlnDebug("SRHQUIC Version Not adapted: tempMsgReceived", tempMsgReceived)
			//lib.PrintlnDebug("SRHQUIC Version Not adapted: tempMsgReceived.QUICStream", tempMsgReceived.QUICStream)
			if tempMsgReceived.QUICStream == nil { // TODO dcruzb: Change to Protocol.Client
				*reset = true
				return
			}
			if isNewConnection, miopPacket := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
				lib.PrintlnDebug("SRHQUIC Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
				*reset = true
				return
			}
			msg.ToAddr = tempMsgReceived.ToAddress //Chn.RemoteAddr().String()
		}
	default:
		{
			*reset = true
			return
		}
	}

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHQUIC Version Not adapted")
	return
}

func (s SRHQUIC) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHQUIC Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr)
	client := srhInfo.Protocol.GetClientFromAddr(msg.ToAddr)
	if client == nil {
		*reset = true
		return
	}
	lib.PrintlnDebug("SRHQUIC Version Not adapted   >>>>> QUIC => msg.ToAddr:", msg.ToAddr, "QUIC Client:", client) //, "AdaptId:", client.AdaptId) // TODO dcruzb: verify impact of removing AdaptId
	msgTemp := msg.Payload.([]byte)

	err := client.Send(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHQUIC Version Not adapted")
	return
}

func (s SRHQUIC) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHQUIC Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	client := srhInfo.Protocol.GetClient(connectionIndex)
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHQUIC Version Not adapted")

		msg, err := client.Receive()
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msg); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHQUIC Version Not adapted")
				break
			}
		}
		if isNewConnection, _ := s.isNewConnection(msg); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			lib.PrintlnDebug("QUIC Is New Connection")
			//miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			//msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			lib.PrintlnDebug("QUIC Before send")
			//s.send(conn, addr, msgPayload)
			lib.PrintlnDebug("QUIC After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		if !*executeForever {
			break
		}
		rcvMessage := messages.ReceivedMessages{Msg: msg, QUICStream: srhInfo.Protocol.GetClient(connectionIndex).Connection().(quic.Stream), ToAddress: srhInfo.Protocol.GetClient(connectionIndex).Address()}
		lib.PrintlnDebug("SRHQUIC Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")

		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHQUIC Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHQUIC Version Not adapted")
}

func (s SRHQUIC) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHQUIC Version Not adapted")
	miop, err := Jsonmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHQUIC) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHQUIC Version Not adapted")
	miop, err := Jsonmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
