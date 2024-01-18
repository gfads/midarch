package middleware

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/gmidarch/development/protocols"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: SRHTLS
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHTLS struct{}

func (s SRHTLS) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	if srhInfo.Protocol == nil {
		srhInfo.Protocol = &protocols.TLS{}
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
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHTLS Version 2 adapted - No connection available")
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
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")
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
			lib.PrintlnDebug("SRHTLS Version Not adapted: tempMsgReceived", tempMsgReceived)
			lib.PrintlnDebug("SRHTLS Version Not adapted: tempMsgReceived.Conn", tempMsgReceived.Conn)
			if tempMsgReceived.Conn == nil { // TODO dcruzb: Change to Protocol.Client
				*reset = true
				return
			}
			if isNewConnection, miopPacket := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
				lib.PrintlnDebug("SRHTLS Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
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

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr)
	client := srhInfo.Protocol.GetClientFromAddr(msg.ToAddr)
	if client == nil {
		*reset = true
		return
	}
	lib.PrintlnDebug("SRHTLS Version Not adapted   >>>>> TLS => msg.ToAddr:", msg.ToAddr, "TLS Client:", client) //, "AdaptId:", client.AdaptId) // TODO dcruzb: verify impact of removing AdaptId
	msgTemp := msg.Payload.([]byte)

	err := client.Send(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	client := srhInfo.Protocol.GetClient(connectionIndex)
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHTLS Version Not adapted")

		msg, err := client.Receive()
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msg); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHTLS Version Not adapted")
				break
			}
		}
		if isNewConnection, _ := s.isNewConnection(msg); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			lib.PrintlnDebug("TLS Is New Connection")
			//miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			//msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			lib.PrintlnDebug("TLS Before send")
			//s.send(conn, addr, msgPayload)
			lib.PrintlnDebug("TLS After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		rcvMessage := messages.ReceivedMessages{Msg: msg, Conn: srhInfo.Protocol.GetClient(connectionIndex).Connection().(net.Conn), ToAddress: srhInfo.Protocol.GetClient(connectionIndex).Address()}
		lib.PrintlnDebug("SRHTLS Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHTLS Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
}

func (s SRHTLS) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHTLS) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
