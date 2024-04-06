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
)

// @Type: SRHRPC
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHRPC struct {
	// Graph exec.ExecGraph
}

func (s SRHRPC) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	if srhInfo.Protocol == nil {
		srhInfo.Protocol = &protocols.RPC{}
		srhInfo.Protocol.StartServer(srhInfo.EndPoint.Host, srhInfo.EndPoint.Port, 2) //shared.MAX_NUMBER_OF_CONNECTIONS)
		//lib.PrintlnInfo("SRHRPC Server Started")
	}

	// // check if a listener has already been created
	// if srhInfo.Ln == nil { // no listen created
	// 	servAddr, err := net.ResolveTCPAddr("tcp", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
	// 	if err != nil {read
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// 	srhInfo.Ln, err = net.ListenTCP("tcp", servAddr)
	// 	if err != nil {
	// 		shared.ErrorHandler(shared.GetFunction(), err.Error())
	// 	}
	// }

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := srhInfo.Protocol.AvailableConnectionFromPool()
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHRPC Version 2 adapted - No connection available")
		time.Sleep(1 * time.Millisecond)
		return
	}

	// RPC is already executed concurrently, dont need go func
	// go func() {
	// lib.PrintlnInfo("SRHRPC Clients Index", availableConenctionIndex)

	client := srhInfo.Protocol.WaitForConnection(availableConenctionIndex)
	// lib.PrintlnInfo("SRHRPC Client", client)

	// Update info
	*info = srhInfo

	// Start goroutine
	if client != nil {
		go s.handler(info, availableConenctionIndex)
	}
	// }()
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

func (s SRHRPC) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")
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
			// lib.PrintlnInfo("SRHRPC Version Not adapted: tempMsgReceived", tempMsgReceived)
			if isNewConnection, _ := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
				// lib.PrintlnInfo("SRHRPC Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
				*reset = true
				return
			}
			// lib.PrintlnInfo("SRHRPC tempMsgReceived.ToAddress", tempMsgReceived.ToAddress)
			msg.ToAddr = tempMsgReceived.ToAddress //Chn.RemoteAddr().String()
		}
	default:
		{
			*reset = true
			return
		}
	}

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

func (s SRHRPC) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr)
	client := srhInfo.Protocol.GetClientFromAddr(msg.ToAddr)
	if client == nil {
		*reset = true
		return
	}
	lib.PrintlnDebug("SRHRPC Version Not adapted   >>>>> RPC => msg.ToAddr:", msg.ToAddr, "RPC Client:", client) //, "AdaptId:", client.AdaptId) // TODO dcruzb: verify impact of removing AdaptId
	msgTemp := msg.Payload.([]byte)

	err := client.Send(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

func (s SRHRPC) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	client := srhInfo.Protocol.GetClient(connectionIndex)
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "FOR", "SRHRPC Version Not adapted")

		msg, err := client.Receive()
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		}
		// lib.PrintlnInfo("SRHRPC got message")
		if changeProtocol, miopPacket := s.isAdapt(msg); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				// lib.PrintlnInfo("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHRPC Version Not adapted")
				break
			}
		}
		if isNewConnection, _ := s.isNewConnection(msg); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			// lib.PrintlnInfo("RPC Is New Connection")
			//miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			//msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			// lib.PrintlnInfo("RPC Before send")
			//s.send(conn, addr, msgPayload)
			// lib.PrintlnInfo("RPC After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		if !*executeForever {
			break
		}
		rcvMessage := messages.ReceivedMessages{Msg: msg, Conn: nil, ToAddress: srhInfo.Protocol.GetClient(connectionIndex).Address()}
		// lib.PrintlnInfo("SRHRPC Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")

		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHRPC Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
}

func (s SRHRPC) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop, err := Jsonmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHRPC) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop, err := Jsonmarshaller{}.Unmarshall(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
		return false, miop
	}
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
