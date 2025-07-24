package middleware

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/components/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: SRHUDP
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHUDP struct{}

func (s SRHUDP) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	if srhInfo.Protocol == nil {
		srhInfo.Protocol = &protocols.UDP{}
		srhInfo.Protocol.StartServer(srhInfo.EndPoint.Host, srhInfo.EndPoint.Port, 1) //shared.MAX_NUMBER_OF_CONNECTIONS)
	}

	// // check if a listener has already been created
	// if srhInfo.Ln == nil { // no listen created
	// 	servAddr, err := net.ResolveUDPAddr("udp", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
	// 	if err != nil {
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// 	srhInfo.Ln, err = net.ListenUDP("udp", servAddr)
	// 	if err != nil {
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// }

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := srhInfo.Protocol.AvailableConnectionFromPool()
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHUDP Version 2 adapted - No connection available")
		time.Sleep(1 * time.Millisecond)
		return
	}

	// go func() {
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

	client := srhInfo.Protocol.WaitForConnection(availableConenctionIndex)
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

	// Update info
	*info = srhInfo

	// Start goroutine
	if client != nil {
		go s.handler(info, availableConenctionIndex)
	}
	// }()
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")
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
			// lib.PrintlnDebug("SRHUDP Version Not adapted: tempMsgReceived", tempMsgReceived)
			// lib.PrintlnDebug("SRHUDP Version Not adapted: tempMsgReceived.Conn", tempMsgReceived.Conn)
			if tempMsgReceived.Conn == nil { // TODO dcruzb: Change to Protocol.Client
				*reset = true
				return
			}
			if isNewConnection, _ := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
				// lib.PrintlnDebug("SRHUDP Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
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

	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// lib.PrintlnDebug("msg.ToAddr", msg.ToAddr)
	client := srhInfo.Protocol.GetClientFromAddr(msg.ToAddr)
	if client == nil {
		*reset = true
		return
	}
	// lib.PrintlnDebug("SRHUDP Version Not adapted   >>>>> UDP => msg.ToAddr:", msg.ToAddr, "UDP Client:", client) //, "AdaptId:", client.AdaptId) // TODO dcruzb: verify impact of removing AdaptId
	msgTemp := msg.Payload.([]byte)

	err := client.Send(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) handler(info *interface{}, connectionIndex int) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	client := srhInfo.Protocol.GetClient(connectionIndex)
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHUDP Version Not adapted")

		msg, err := client.Receive()
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		}
		// lib.PrintlnDebug("Message received")

		if changeProtocol, miopPacket := s.isAdapt(msg); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
				break
			}
		}
		if isNewConnection, miopPacket := s.isNewConnection(msg); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			// lib.PrintlnDebug("UDP Is New Connection")
			miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			msgPayload := middleware.Gobmarshaller{}.Marshall(miopPacket)
			// lib.PrintlnDebug("UDP Before send")
			err := client.Send(msgPayload)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			// lib.PrintlnDebug("UDP After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		if !*executeForever {
			break
		}
		rcvMessage := messages.ReceivedMessages{Msg: msg, Conn: srhInfo.Protocol.GetClient(connectionIndex).Connection().(net.Conn), ToAddress: srhInfo.Protocol.GetClient(connectionIndex).Address()}
		// lib.PrintlnDebug("SRHUDP Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")

		srhInfo.RcvedMessages <- rcvMessage
		// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHUDP Version Not adapted")
	}
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
}

func (s SRHUDP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	miop, err := middleware.Gobmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHUDP) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	miop, err := middleware.Gobmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
