package srhtcp

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/gfads/midarch/pkg/gmidarch/development/components/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
)

// @Type: SRHTCP
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHTCP struct{}

func (s SRHTCP) availableConnectionFromPool(clientsPtr *[]*messages.Client, ip string) (bool, int) {
	clients := *clientsPtr

	if ip != "" {
		for idx, client := range clients {
			if client.Ip == ip {
				return true, idx
			}
			if client.UDPConnection != nil {
				client.UDPConnection = nil
				return true, idx
			}
		}
	}

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < 10 { //shared.MAX_NUMBER_OF_CONNECTIONS { TODO: dcruzb go back the env var
		client := messages.Client{
			Ip:            "",
			Connection:    nil,
			UDPConnection: nil,
		}
		*clientsPtr = append(clients, &client)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(*clientsPtr) - 1
	}

	for idx, client := range clients {
		if client == nil {
			client := messages.Client{
				Ip:            "",
				Connection:    nil,
				UDPConnection: nil,
			}
			clients[idx] = &client
			return true, idx
		}
		if client.UDPConnection != nil {
			client.UDPConnection.Close()
			client.UDPConnection = nil
			return true, idx
		}
	}

	return false, -1
}

func (s SRHTCP) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

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

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, "")
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//fmt.Println("------------------------------>", shared.GetFunction(), "end", "SRHTCP Version 2 adapted - No connection available")
		return
	}

	go func() {
		//fmt.Println("SRHTCP Version 2 adapted >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex, " ip:", srhInfo.Clients[availableConenctionIndex].Ip)

		// Accept connections
		conn, err := srhInfo.Ln.Accept()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		srhInfo.Conns = append(srhInfo.Conns, conn)
		//srhInfo.CurrentConn = conn

		//fmt.Println("SRHTCP Version 2 adapted >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Buscou nova conexão, ip:", conn.RemoteAddr().String())
		//connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, conn.RemoteAddr().String())
		////log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(
		////srhInfo.Clients))
		//if !connectionAvailable {
		//	fmt.Println("------------------------------>", shared.GetFunction(), "end", "SRHTCP Version 2 adapted - No connection available")
		//	return
		//}

		client := srhInfo.Clients[availableConenctionIndex]
		client.Ip = conn.RemoteAddr().String()
		client.Connection = conn
		//fmt.Println("SRHTCP Version 2 adapted >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

		// Update info
		*info = srhInfo

		// Start goroutine
		go s.handler(info, availableConenctionIndex)
	}()
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version 2 adapted")
	return
}

func (s SRHTCP) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version 2 adapted")
	//fmt.Println(shared.GetFunction(), "HERE")
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
			//fmt.Println("SRHTCP Version 2 adapted: tempMsgReceived", tempMsgReceived)
			//fmt.Println("SRHTCP Version 2 adapted: tempMsgReceived.Conn", tempMsgReceived.Conn)
			if tempMsgReceived.Conn == nil {
				*reset = true
				return
			}
			msg.ToAddr = tempMsgReceived.ToAddress //Conn.RemoteAddr().String()
		}
	default:
		{
			*reset = true
			return
		}
	}

	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version 2 adapted")
	return
}

func (s SRHTCP) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	//fmt.Println("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	client := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients)
	conn := client.Connection //srhInfo.CurrentConn
	if conn == nil {
		*reset = true
		return
	}
	//fmt.Println("SRHTCP Version 2 adapted   >>>>> TCP => msg.ToAddr:", msg.ToAddr, "TCP conn:", conn, "AdaptId:", client.AdaptId)
	msgTemp := msg.Payload.([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//json := middleware.Gobmarshaller{}
	//unmarshalledMsg := json.Unmarshall(msgTemp)
	//fmt.Println("<<<<<<<<<<<<  <<<<<<<<<<  <<<<<<<<<  SRHTCP Version 2 adapted => Msg: ", unmarshalledMsg.Bd.RepBody.OperationResult)
	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version 2 adapted")
	return
}

func (s SRHTCP) handler(info *interface{}, connectionIndex int) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version 2 adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR", "SRHTCP Version 2 adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := conn.Read(size)
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "Read efetuado", "SRHTCP Version 2 adapted")
		if err == io.EOF {
			srhInfo.Clients[connectionIndex] = nil
			//fmt.Println("Não Vai matar o app EOF")
			break
		} else if err != nil && err != io.EOF {
			//fmt.Println("Vai matar o app, erro mas não EOF")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = conn.Read(msgTemp)
		if err == io.EOF {
			srhInfo.Clients[connectionIndex] = nil
			//fmt.Println("Não Vai matar o app EOF")
			break
		} else if err != nil && err != io.EOF {
			//fmt.Println("Vai matar o app, erro mas não EOF")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		if changeProtocol, miopPacket := s.isAdapt(msgTemp); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				//fmt.Println("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHTCP Version 2 adapted")
				break
			}
		}

		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Conn: conn, ToAddress: srhInfo.Clients[connectionIndex].Ip}
		//fmt.Println("SRHTCP Version 2 adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHTCP Version 2 adapted")
	}
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version 2 adapted")
}

func (s SRHTCP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop, _ := middleware.Gobmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}
